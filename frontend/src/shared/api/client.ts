// Generated from backend/openapi/openapi.yaml.
// Manual snapshot for current stage with stub pages.

export const apiPaths = {
  adminLogin: "/api/admin/auth/login",
  adminTenants: "/api/admin/tenants",
  adminTenantById: (id: string) => `/api/admin/tenants/${id}`,
  adminStats: "/api/admin/stats",
  clientLogin: "/api/client/auth/login",
  clientRegister: "/api/client/auth/register",
  products: "/api/client/products",
  productById: (id: string) => `/api/client/products/${id}`,
  variants: "/api/client/variants",
  variantById: (id: string) => `/api/client/variants/${id}`,
  variantComponents: (id: string) => `/api/client/variants/${id}/components`,
  variantStock: (id: string) => `/api/client/variants/${id}/stock`,
  warehouses: "/api/client/warehouses",
  warehouseById: (id: string) => `/api/client/warehouses/${id}`,
  documents: "/api/client/documents",
  documentById: (id: string) => `/api/client/documents/${id}`,
  documentPost: (id: string) => `/api/client/documents/${id}/post`,
  documentCancel: (id: string) => `/api/client/documents/${id}/cancel`,
  reportProfit: "/api/client/reports/profit",
  reportStock: "/api/client/reports/stock",
  reportTopProducts: "/api/client/reports/top-products",
  reportMovements: "/api/client/reports/movements",
};
