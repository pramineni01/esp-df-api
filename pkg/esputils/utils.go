package esputils

import "context"

func GetUserIDFromContext(ctx context.Context) (string, error) {
	// TODO: implement retrieval of UserID from keycloak token in ctx
	return "test", nil
}

func GetAccessViewIDFromContext(ctx context.Context) (int, error) {
	// TODO: implement retrieval of access scope from keycloak token in ctx
	return 0, nil
}
