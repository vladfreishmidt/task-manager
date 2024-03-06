package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const userID = contextKey("userID")
