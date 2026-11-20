import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function connectSeller() {
  const { ticket, launch_url } = await postJSON(apiPaths.allegroConnect, {});
  const form = document.createElement("form");
  const field = document.createElement("input");
  form.method = "post";
  form.action = launch_url;
  field.type = "hidden";
  field.name = "ticket";
  field.value = ticket;
  form.append(field);
  document.body.append(form);
  form.submit();
  form.remove();
}
