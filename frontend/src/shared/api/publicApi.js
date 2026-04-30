const apiScopes = Object.freeze({
  admin: "/api/admin",
  client: "/api/client",
});

export const apiPaths = Object.freeze({
  adminLogin: `${apiScopes.admin}/auth/login`,
  adminTenants: `${apiScopes.admin}/tenants`,
  adminTenantById: (id) => `${apiScopes.admin}/tenants/${id}`,
  adminStats: `${apiScopes.admin}/stats`,
  clientLogin: `${apiScopes.client}/auth/login`,
  clientRegister: `${apiScopes.client}/auth/register`,
  currencies: `${apiScopes.client}/currencies`,
  currencyOptions: `${apiScopes.client}/currency-options`,
  products: `${apiScopes.client}/products`,
  productById: (id) => `${apiScopes.client}/products/${id}`,
  bundles: `${apiScopes.client}/bundles`,
  bundleById: (id) => `${apiScopes.client}/bundles/${id}`,
  bundleComponents: (id) => `${apiScopes.client}/bundles/${id}/components`,
  variants: `${apiScopes.client}/variants`,
  variantById: (id) => `${apiScopes.client}/variants/${id}`,
  variantStock: (id) => `${apiScopes.client}/variants/${id}/stock`,
  warehouses: `${apiScopes.client}/warehouses`,
  warehouseById: (id) => `${apiScopes.client}/warehouses/${id}`,
  customers: `${apiScopes.client}/customers`,
  customerById: (id) => `${apiScopes.client}/customers/${id}`,
  documents: `${apiScopes.client}/documents`,
  documentById: (id) => `${apiScopes.client}/documents/${id}`,
  documentPost: (id) => `${apiScopes.client}/documents/${id}/post`,
  documentCancel: (id) => `${apiScopes.client}/documents/${id}/cancel`,
  reportProfit: `${apiScopes.client}/reports/profit`,
  reportStock: `${apiScopes.client}/reports/stock`,
  reportTopProducts: `${apiScopes.client}/reports/top-products`,
  reportMovements: `${apiScopes.client}/reports/movements`,
});

export function assertPublicApiPath(path) {
  if (isPublicApiPath(path)) {
    return;
  }
  throw new Error(`Frontend can call only public API paths: ${path}`);
}

export function isAdminApiPath(path) {
  return startsWithScope(path, apiScopes.admin);
}

export function isClientApiPath(path) {
  return startsWithScope(path, apiScopes.client);
}

function isPublicApiPath(path) {
  return isAdminApiPath(path) || isClientApiPath(path);
}

function startsWithScope(path, scope) {
  return typeof path === "string" && path.startsWith(scope);
}
