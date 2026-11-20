import { useMemo } from "react";
import { apiPaths } from "@smarterp/api/publicApi";
import useApiResource from "@smarterp/api/useApiResource";

export function useCountryOptions(locale) {
  const result = useApiResource(apiPaths.countryOptions);
  const options = useMemo(() => {
    const names = new Intl.DisplayNames([locale], { type: "region" });
    return (result.data ?? []).map((code) => [code, names.of(code)]).sort((a, b) => a[1].localeCompare(b[1], locale));
  }, [locale, result.data]);
  return { ...result, options };
}
