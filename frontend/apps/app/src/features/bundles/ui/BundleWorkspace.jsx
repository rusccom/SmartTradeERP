import SmartBundleSettings from "../../integrations/smartbundle/SmartBundleSettings";
import BundleComponentsEditor from "./BundleComponentsEditor";

export default function BundleWorkspace(props) {
  return (
    <div className="bundle-workspace">
      <BundleComponentsEditor {...props} />
      {props.bundle && <SmartBundleSettings key={props.bundle.variant_id} bundle={props.bundle} />}
    </div>
  );
}
