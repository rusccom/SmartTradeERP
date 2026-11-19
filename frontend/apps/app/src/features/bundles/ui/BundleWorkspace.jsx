import { usePlugins } from "@smarterp/model/plugins/usePlugins";
import SmartBundleSettings from "../../integrations/smartbundle/SmartBundleSettings";
import BundleComponentsEditor from "./BundleComponentsEditor";

export default function BundleWorkspace(props) {
  const { access } = usePlugins();
  const configured = access?.plugins.some(plugin => plugin.key === "smartbundle" && plugin.enabled);
  return (
    <div className="bundle-workspace">
      <BundleComponentsEditor {...props} />
      {configured && props.bundle && <SmartBundleSettings key={props.bundle.variant_id} bundle={props.bundle} />}
    </div>
  );
}
