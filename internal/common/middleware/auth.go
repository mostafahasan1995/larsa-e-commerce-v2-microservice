// Package middleware provides authentication and authorization middleware for HTTP handlers.
// It implements both basic authentication and role-based access control (RBAC) mechanisms
// for securing API endpoints in a microservices architecture.
package middleware

import (
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"net/http"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth creates a middleware that performs basic authentication.
// It takes a variable number of restriction strings that specify required permissions.
//
// The middleware:
// 1. Extracts authentication headers from the request
// 2. Validates credentials asynchronously
// 3. Creates a user context with authenticated user information
// 4. Passes the authenticated request to the next handler
//
// If authentication fails, it returns HTTP 403 Forbidden status.
//
// Example usage:
//
//	http.Handle("/api/resource", Auth("read")(resourceHandler))
func Auth(restrictions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerParams := common.ExtractHeaderParams(r)
			routeRestrictions := make([]string, len(restrictions))
			copy(routeRestrictions, restrictions)

			if l := len(routeRestrictions); l == 0 {
				if headerParams.AccessToken == "" {
					next.ServeHTTP(w, r)
					return
				}
				routeRestrictions = append(routeRestrictions, "authenticate")
			}

			c := make(chan common.Credentials)

			go common.Guard(c, headerParams, routeRestrictions)

			credentials := <-c

			if credentials.Err != nil {
				http.Error(w, "not authenticated", http.StatusForbidden)
				return
			}

			_user := credentials.User

			userId, err := primitive.ObjectIDFromHex(_user.Id)
			if err != nil {
				http.Error(w, "error get user id", http.StatusForbidden)
				return
			}

			user := &types.User{
				Id:       userId,
				UserData: _user,
				Type:     enums.Normal,
			}

			ctx := util.SetReqUser(r.Context(), user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware is a type that defines the signature for HTTP middleware functions.
// It takes an http.Handler and returns a wrapped http.Handler with additional functionality.
type Middleware func(http.Handler) http.Handler

// Restriction defines the capability requirements for different user roles.
// Each field specifies the permission string required for that role level.
type Restriction struct {
	SuperAdminCapability      string // Permission required for SuperAdmin access
	OwnerCapability           string // Permission required for Owner access
	MemberCapability          string // Permission required for Member access
	BuyerCapability           string // Permission required for Buyer access
	VendorCapability          string // Permission required for Vendor access
	AgentCapability           string // Permission required for Agent access
	ShippingCompanyCapability string // Permission required for ShippingCompany access
}

// AdminAuth creates a middleware that implements role-based authentication.
// It performs parallel authentication checks for SuperAdmin, Owner, and Member roles,
// assigning the highest successful privilege level to the user.
//
// Parameters:
//   - restriction: Struct containing required capabilities for each role level
//
// The middleware:
// 1. Extracts header parameters from the request
// 2. Performs concurrent authentication checks for all roles
// 3. Determines the highest privilege level the user has access to
// 4. Adds user information and role type to the request context
//
// If all authentication checks fail, returns HTTP 403 Forbidden status.
//
// Example usage:
//
//	restrictions := Restriction{
//	    SuperAdminCapability: "manage_all",
//	    OwnerCapability:     "manage_own",
//	    MemberCapability:    "view"
//	}
//	http.Handle("/api/admin/resource", AdminAuth(restrictions)(adminResourceHandler))
func AdminAuth(restriction Restriction) Middleware {
	if configuration.EnvironmentVariables.Mode == "debug" {
		restriction.SuperAdminCapability = "authenticate"

		restriction.OwnerCapability = "authenticate"
		restriction.MemberCapability = "authenticate"

		restriction.BuyerCapability = "authenticate"
		restriction.VendorCapability = "authenticate"

		restriction.AgentCapability = "authenticate"
		restriction.ShippingCompanyCapability = "authenticate"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerParams := common.ExtractHeaderParams(r)
			channelUser := make(chan common.Credentials)
			go common.Guard(channelUser, headerParams, []string{"authenticate"})
			userCredentials := <-channelUser
			if userCredentials.Err != nil {
				http.Error(w, "not Authenticated", http.StatusUnauthorized)
				return
			}

			channelAdmin := make(chan common.Credentials)
			go common.Guard(channelAdmin, headerParams, []string{"authenticate", restriction.SuperAdminCapability})
			channelOwner := make(chan common.Credentials)

			go common.Guard(channelOwner, headerParams, []string{"authenticate", restriction.OwnerCapability})

			channelMember := make(chan common.Credentials)
			go common.Guard(channelMember, headerParams, []string{"authenticate", restriction.MemberCapability})

			channelBuyer := make(chan common.Credentials)
			go common.Guard(channelBuyer, headerParams, []string{"authenticate", restriction.BuyerCapability})

			channelVendor := make(chan common.Credentials)
			go common.Guard(channelVendor, headerParams, []string{"authenticate", restriction.VendorCapability})

			channelAgent := make(chan common.Credentials)
			go common.Guard(channelAgent, headerParams, []string{"authenticate", restriction.AgentCapability})

			channelShippingCompany := make(chan common.Credentials)
			go common.Guard(channelShippingCompany, headerParams, []string{"authenticate", restriction.ShippingCompanyCapability})

			adminCredentials := <-channelAdmin
			ownerCredentials := <-channelOwner
			memberCredentials := <-channelMember
			buyerCredentials := <-channelBuyer
			vendorCredentials := <-channelVendor
			agentCredentials := <-channelAgent
			shippingCompanyCredentials := <-channelShippingCompany
			// fmt.Printf("AccessToken: %v\n", headerParams.AccessToken)
			// fmt.Printf("Admin Credentials: %v, UserID: %v, Error: %v\n", adminCredentials.Authentication, adminCredentials.User.Id, adminCredentials.Err)
			// fmt.Printf("Owner Credentials: %v, UserID: %v, Error: %v\n", ownerCredentials.Authentication, ownerCredentials.User.Id, ownerCredentials.Err)
			// fmt.Printf("Member Credentials: %v, UserID: %v, Error: %v\n", memberCredentials.Authe	ntication, memberCredentials.User.Id, memberCredentials.Err)
			// fmt.Printf("Buyer Credentials: %v, UserID: %v, Error: %v\n", buyerCredentials.Authentication, buyerCredentials.User.Id, buyerCredentials.Err)
			// fmt.Printf("Vendor Credentials: %v, UserID: %v, Error: %v\n", vendorCredentials.Authentication, vendorCredentials.User.Id, vendorCredentials.Err)

			if adminCredentials.Err != nil && ownerCredentials.Err != nil && memberCredentials.Err != nil && buyerCredentials.Err != nil && vendorCredentials.Err != nil && agentCredentials.Err != nil && shippingCompanyCredentials.Err != nil {
				http.Error(w, "not Authorized", http.StatusForbidden)
				return
			}

			if !adminCredentials.Authentication && !ownerCredentials.Authentication && !memberCredentials.Authentication && !buyerCredentials.Authentication && !vendorCredentials.Authentication && !agentCredentials.Authentication && !shippingCompanyCredentials.Authentication {
				http.Error(w, "not Authorized", http.StatusForbidden)
				return
			}
			var userType enums.UserType

			var user common.User
			if adminCredentials.Authentication {
				user = adminCredentials.User
				userType = enums.SuperAdmin
			} else if vendorCredentials.Authentication {
				user = vendorCredentials.User
				userType = enums.Vendor
			} else if agentCredentials.Authentication {
				user = agentCredentials.User
				userType = enums.Agent
			} else if shippingCompanyCredentials.Authentication {
				user = shippingCompanyCredentials.User
				userType = enums.ShippingCompany
			} else if memberCredentials.Authentication {
				user = memberCredentials.User
				userType = enums.Member
			} else if buyerCredentials.Authentication {
				user = buyerCredentials.User
				userType = enums.Normal
			}

			userId, err := primitive.ObjectIDFromHex(user.Id)
			if err != nil {
				http.Error(w, "error get user id", http.StatusForbidden)
				return
			}

			ctx := util.SetReqUser(r.Context(), &types.User{
				Id:       userId,
				UserData: user,
				Type:     userType,
			})
			// ctx = context.WithValue(r.Context(), common.CommonContextKey("user"), user)
			// ctx = context.WithValue(ctx, common.CommonContextKey("userType"), userType)
			// next.ServeHTTP(w, r.Clone(ctx))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
