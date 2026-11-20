import { useCountryOptions } from "@smarterp/model/countryOptions";
import { apiPaths } from "@smarterp/api/publicApi";
import useApiResource from "@smarterp/api/useApiResource";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function useVatCountries() {
  const { locale } = useI18n();
  const result = useApiResource(apiPaths.vatCountries, true);
  const dictionary = useCountryOptions(locale);
  const enabled = dictionary.options.filter(([code]) => result.data?.includes(code));
  return { data: result.data, options: dictionary.options, enabled, loading: result.loading || dictionary.loading,
    error: result.error || dictionary.error, retry: () => { result.retry(); dictionary.retry(); } };
}
