package auth_v1_test

import (
	"auth-service/internal/api/auth_v1"
	"auth-service/internal/config"
	"auth-service/internal/model"
	service_mocks "auth-service/internal/service/mocks"
	pb "auth-service/pkg/api/auth_v1"
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func isValidToken(tokenString, yourSecret string) bool {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(yourSecret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}
	return false
}

func createClient(t *testing.T) pb.AuthClient {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	sau := auth_v1.AuthServer{}
	pb.RegisterAuthServer(srv, &sau)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		log.Fatalf("grpc.DialContext %v", err)
	}

	client := pb.NewAuthClient(conn)

	return client
}

func TestAuthServer_Login(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockUserService, email, password string)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	client := createClient(t)

	cases := []struct {
		name    string
		args    *pb.AuthRequest
		want    bool
		wantErr bool
	}{
		{
			name:    "should return JWT valid token",
			args:    &pb.AuthRequest{Email: "test@test.org", Password: "password123"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "should return error when email is empty",
			args:    &pb.AuthRequest{Email: "", Password: "password123"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "should return error when password is empty",
			args:    &pb.AuthRequest{Email: "test@test.org", Password: ""},
			want:    false,
			wantErr: true,
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			res, err := client.Login(context.Background(), item.args)
			if item.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, item.want, isValidToken(res.Token, cfg.App.SecretKey))
			}
		})
	}
}

func TestAuthServer_Register(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockUserService, user *model.User)

	client := createClient(t)

	cases := []struct {
		name      string
		args      *pb.RegisterRequest
		want      *pb.RegisterResponse
		wantErr   bool
		phraseErr string
	}{
		{
			name: "should return user object",
			args: &pb.RegisterRequest{
				User: &pb.User{
					Email:    "test@test.org",
					Password: "test123",
					Username: "test",
					Roles:    "admin, test",
				},
			},
			want: &pb.RegisterResponse{
				User: &pb.User{
					Email:    "test@test.org",
					Username: "test",
					Roles:    "admin, test",
				},
			},
			wantErr:   false,
			phraseErr: "",
		},
		{
			name: "should return error when email is empty",
			args: &pb.RegisterRequest{User: &pb.User{
				Email:    "",
				Password: "test123",
				Username: "test",
				Roles:    "admin, test",
			}},
			want:      &pb.RegisterResponse{},
			wantErr:   true,
			phraseErr: "email is required",
		},
		{
			name: "should return error when password is empty",
			args: &pb.RegisterRequest{User: &pb.User{
				Email:    "test@test.org",
				Password: "",
				Username: "test",
				Roles:    "admin, test",
			}},
			want:      &pb.RegisterResponse{},
			wantErr:   true,
			phraseErr: "password is required",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			res, err := client.Register(context.Background(), item.args)
			if item.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, item.want.User.Email, res.User.Email)
				require.Equal(t, item.want.User.Username, res.User.Username)
				require.Equal(t, item.want.User.Roles, res.User.Roles)
				require.NotEmpty(t, res.User.Uuid)
				require.NotEmpty(t, res.User.CreatedAt)
				require.NotEmpty(t, res.User.UpdatedAt)
			}
		})
	}
}
