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
  orders: `${apiScopes.client}/orders`,
  order: (id) => `${apiScopes.client}/orders/${encodeURIComponent(id)}`,
  orderQuote: (id) => id ? `${apiScopes.client}/orders/${encodeURIComponent(id)}/quote` : `${apiScopes.client}/orders/quote`,
  orderPayment: (id) => `${apiScopes.client}/orders/${encodeURIComponent(id)}/payments`,
  invoiceFormats: `${apiScopes.client}/orders/invoice-formats`,
  orderInvoice: (id) => `${apiScopes.client}/orders/${encodeURIComponent(id)}/invoice`,
  polishSeller: `${apiScopes.client}/regions/pl/settings`,
  invoiceDelivery: (id) => `${apiScopes.client}/regions/pl/invoices/${encodeURIComponent(id)}/ksef`,
  checkInvoiceDelivery: (id) => `${apiScopes.client}/regions/pl/invoices/${encodeURIComponent(id)}/ksef/check`,
  currencyBase: `${apiScopes.client}/currencies/base`,
  countryOptions: `${apiScopes.client}/country-options`,
  currencyOptions: `${apiScopes.client}/currency-options`,
  clientPlugins: `${apiScopes.client}/plugins`,
  inpostSettings: `${apiScopes.client}/integrations/inpost`,
  allegro: `${apiScopes.client}/integrations/allegro`,
  allegroToken: `${apiScopes.client}/integrations/allegro/token`,
  allegroPrivateToken: `${apiScopes.client}/integrations/allegro/private/refresh-token`,
  allegroPrivateDevice: `${apiScopes.client}/integrations/allegro/private/device`,
  allegroPrivatePoll: `${apiScopes.client}/integrations/allegro/private/device/check`,
  allegroConnection: (id) => `${apiScopes.client}/integrations/allegro/connections/${encodeURIComponent(id)}`,
  allegroConnect: `${apiScopes.client}/integrations/allegro/oauth/connect`,
  allegroCheck: `${apiScopes.client}/integrations/allegro/check`,
  clientPlugin: (key) => `${apiScopes.client}/plugins/${encodeURIComponent(key)}`,
  adminTenantPlugins: (id) => `${apiScopes.admin}/tenants/${encodeURIComponent(id)}/plugins`,
  adminTenantPlugin: (id, key) => `${apiScopes.admin}/tenants/${encodeURIComponent(id)}/plugins/${encodeURIComponent(key)}`,
  storefrontSettings: `${apiScopes.client}/storefront/settings`,
  storefrontThemes: `${apiScopes.client}/storefront/themes`,
  storefrontDraft: `${apiScopes.client}/storefront/draft`,
  storefrontPublish: `${apiScopes.client}/storefront/publish`,
  storefrontPreview: `${apiScopes.client}/storefront/preview`,
  vatCountries: `${apiScopes.client}/catalog/vat-countries`,
  vatCountry: (code) => `${apiScopes.client}/catalog/vat-countries/${encodeURIComponent(code)}`,
  products: `${apiScopes.client}/products`,
  productById: (id) => `${apiScopes.client}/products/${id}`,
  productMedia: (id) => `${apiScopes.client}/products/${id}/media`,
  productMediaComplete: (id, mediaID) => `${apiScopes.client}/products/${id}/media/${mediaID}/complete`,
  productMediaItem: (id, mediaID) => `${apiScopes.client}/products/${id}/media/${mediaID}`,
  productMediaPrimary: (id, mediaID) => `${apiScopes.client}/products/${id}/media/${mediaID}/primary`,
  bundles: `${apiScopes.client}/bundles`,
  bundleById: (id) => `${apiScopes.client}/bundles/${id}`,
  bundleStorefront: (id) => `${apiScopes.client}/bundles/${id}/storefront`,
  bundleComponents: (id) => `${apiScopes.client}/bundles/${id}/components`,
  variants: `${apiScopes.client}/variants`,
  variantById: (id) => `${apiScopes.client}/variants/${id}`,
  variantStock: (id) => `${apiScopes.client}/variants/${id}/stock`,
  variantMedia: (id) => `${apiScopes.client}/variants/${id}/media`,
  variantMediaComplete: (id, mediaID) => `${apiScopes.client}/variants/${id}/media/${mediaID}/complete`,
  variantMediaItem: (id, mediaID) => `${apiScopes.client}/variants/${id}/media/${mediaID}`,
  variantMediaPrimary: (id, mediaID) => `${apiScopes.client}/variants/${id}/media/${mediaID}/primary`,
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
