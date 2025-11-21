<!-- 6b4b29cd-0480-4044-8a16-8a74fbaa7dff c3575063-7064-43c1-8853-ce0d01e9078c -->
# Fix Warehouse Bills 404 Endpoint

## Problem

The frontend is making a request to `GET /api/v1/franchises/1/warehouse-bills` but receiving a 404 response. The route is registered in `router.go` at line 199, but the server may not have the updated routes loaded.

## Root Cause Analysis

1. The route `franchises.GET("/:franchiseId/warehouse-bills", ...)` is correctly registered in `internal/presentation/http/router/router.go:199`
2. The handler method `ListFranchiseWarehouseBills` exists in `internal/presentation/http/handler/warehousebill_handler.go:148`
3. The service method exists in `internal/application/warehousebill/service.go:659`
4. The backend server likely needs to be restarted to load the new routes

## Solution

### Step 1: Verify Route Registration Order

- Check if route ordering in Gin could cause conflicts
- In Gin, more specific routes should generally come before less specific ones
- Current order: `/:franchiseId` (line 179) comes before `/:franchiseId/warehouse-bills` (line 199)
- This should be fine as Gin matches full paths, but we should verify

### Step 2: Ensure Backend Server Restart

- The backend server must be restarted after adding new routes
- Verify the server is running the latest compiled code
- Check for any compilation errors preventing route registration

### Step 3: Test Route Accessibility

- Test the endpoint directly with curl or Postman
- Verify authentication middleware is not blocking the request
- Check server logs for any route registration errors

## Files to Check/Modify

- `internal/presentation/http/router/router.go` - Verify route registration
- `internal/presentation/http/handler/warehousebill_handler.go` - Verify handler exists
- `cmd/api/main.go` - Verify handler initialization
- Backend server process - Ensure it's running latest code

## Expected Outcome

After restarting the backend server with the updated routes, the endpoint `GET /api/v1/franchises/:franchiseId/warehouse-bills` should return a 200 response with the list of warehouse bills for that franchise.

### To-dos

- [ ] Verify route ordering in router.go - ensure warehouse-bills routes are properly positioned
- [ ] Verify backend server has been restarted with latest code containing warehouse bill routes
- [ ] Test the endpoint GET /api/v1/franchises/1/warehouse-bills to confirm it works after server restart