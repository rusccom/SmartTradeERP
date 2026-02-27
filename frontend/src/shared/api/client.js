// Generated from backend/openapi/openapi.yaml.
// Manual snapshot for current stage with stub pages.

export const apiPaths = {
  adminLogin: "/api/admin/auth/login",
  adminTenants: "/api/admin/tenants",
  adminTenantById: (id) => `/api/admin/tenants/${id}`,
  adminStats: "/api/admin/stats",
  clientLogin: "/api/client/auth/login",
  clientRegister: "/api/client/auth/register",
  products: "/api/client/products",
  productById: (id) => `/api/client/products/${id}`,
  variants: "/api/client/variants",
  variantById: (id) => `/api/client/variants/${id}`,
  variantComponents: (id) => `/api/client/variants/${id}/components`,
  variantStock: (id) => `/api/client/variants/${id}/stock`,
  warehouses: "/api/client/warehouses",
  warehouseById: (id) => `/api/client/warehouses/${id}`,
  customers: "/api/client/customers",
  customerById: (id) => `/api/client/customers/${id}`,
  documents: "/api/client/documents",
  documentById: (id) => `/api/client/documents/${id}`,
  documentPost: (id) => `/api/client/documents/${id}/post`,
  documentCancel: (id) => `/api/client/documents/${id}/cancel`,
  reportProfit: "/api/client/reports/profit",
  reportStock: "/api/client/reports/stock",
  reportTopProducts: "/api/client/reports/top-products",
  reportMovements: "/api/client/reports/movements",
};
