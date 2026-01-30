package main

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 1. Define Public Methods (No Token Needed)
var publicMethods = map[string]bool{
	"/user.UserService/Login":    true,
	"/user.UserService/Register": true,
}

// 2. Define Admin-Only Methods
var adminMethods = map[string]bool{
	"/user.UserService/CreateUser": true,
	"/user.UserService/UpdateUser": true,
	"/user.UserService/DeleteUser": true,
	"/user.UserService/GetUser":    true, // <--- Add this
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// A. Allow Public Methods
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// B. Get Metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata missing")
	}

	// C. Get Token
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "token missing")
	}

	tokenString := values[0]
	if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
		tokenString = tokenString[7:]
	}

	// D. Validate Token & Parse Claims
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// --- NEW: ROLE CHECK ---
	// E. If method requires Admin, check the role
	if adminMethods[info.FullMethod] {
		if strings.ToLower(claims.Role) != "admin" {
			return nil, status.Errorf(codes.PermissionDenied, "Access Denied: You are not an admin")
		}
	}

	// F. Success
	return handler(ctx, req)
}
