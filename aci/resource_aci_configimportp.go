package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciConfigurationImportPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciConfigurationImportPolicyCreate,
		Update: resourceAciConfigurationImportPolicyUpdate,
		Read:   resourceAciConfigurationImportPolicyRead,
		Delete: resourceAciConfigurationImportPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConfigurationImportPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fail_on_decrypt_errors": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"import_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"import_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"snapshot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_config_rs_import_source": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_trig_rs_triggerable": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_config_rs_remote_path": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteConfigurationImportPolicy(client *client.Client, dn string) (*models.ConfigurationImportPolicy, error) {
	configImportPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	configImportP := models.ConfigurationImportPolicyFromContainer(configImportPCont)

	if configImportP.DistinguishedName == "" {
		return nil, fmt.Errorf("ConfigurationImportPolicy %s not found", configImportP.DistinguishedName)
	}

	return configImportP, nil
}

func setConfigurationImportPolicyAttributes(configImportP *models.ConfigurationImportPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(configImportP.DistinguishedName)
	d.Set("description", configImportP.Description)
	configImportPMap, _ := configImportP.ToMap()

	d.Set("name", configImportPMap["name"])

	d.Set("admin_st", configImportPMap["adminSt"])
	d.Set("annotation", configImportPMap["annotation"])
	d.Set("fail_on_decrypt_errors", configImportPMap["failOnDecryptErrors"])
	d.Set("file_name", configImportPMap["fileName"])
	d.Set("import_mode", configImportPMap["importMode"])
	d.Set("import_type", configImportPMap["importType"])
	d.Set("name_alias", configImportPMap["nameAlias"])
	d.Set("snapshot", configImportPMap["snapshot"])
	return d
}

func resourceAciConfigurationImportPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	configImportP, err := getRemoteConfigurationImportPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setConfigurationImportPolicyAttributes(configImportP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConfigurationImportPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ConfigurationImportPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configImportPAttr := models.ConfigurationImportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configImportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configImportPAttr.Annotation = Annotation.(string)
	}
	if FailOnDecryptErrors, ok := d.GetOk("fail_on_decrypt_errors"); ok {
		configImportPAttr.FailOnDecryptErrors = FailOnDecryptErrors.(string)
	}
	if FileName, ok := d.GetOk("file_name"); ok {
		configImportPAttr.FileName = FileName.(string)
	}
	if ImportMode, ok := d.GetOk("import_mode"); ok {
		configImportPAttr.ImportMode = ImportMode.(string)
	}
	if ImportType, ok := d.GetOk("import_type"); ok {
		configImportPAttr.ImportType = ImportType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configImportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configImportPAttr.Snapshot = Snapshot.(string)
	}
	configImportP := models.NewConfigurationImportPolicy(fmt.Sprintf("fabric/configimp-%s", name), "uni", desc, configImportPAttr)

	err := aciClient.Save(configImportP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToconfigRsImportSource, ok := d.GetOk("relation_config_rs_import_source"); ok {
		relationParam := relationToconfigRsImportSource.(string)
		err = aciClient.CreateRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_config_rs_import_source")
		d.Partial(false)

	}
	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_trig_rs_triggerable")
		d.Partial(false)

	}
	if relationToconfigRsRemotePath, ok := d.GetOk("relation_config_rs_remote_path"); ok {
		relationParam := relationToconfigRsRemotePath.(string)
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_config_rs_remote_path")
		d.Partial(false)

	}

	d.SetId(configImportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciConfigurationImportPolicyRead(d, m)
}

func resourceAciConfigurationImportPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ConfigurationImportPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configImportPAttr := models.ConfigurationImportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configImportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configImportPAttr.Annotation = Annotation.(string)
	}
	if FailOnDecryptErrors, ok := d.GetOk("fail_on_decrypt_errors"); ok {
		configImportPAttr.FailOnDecryptErrors = FailOnDecryptErrors.(string)
	}
	if FileName, ok := d.GetOk("file_name"); ok {
		configImportPAttr.FileName = FileName.(string)
	}
	if ImportMode, ok := d.GetOk("import_mode"); ok {
		configImportPAttr.ImportMode = ImportMode.(string)
	}
	if ImportType, ok := d.GetOk("import_type"); ok {
		configImportPAttr.ImportType = ImportType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configImportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configImportPAttr.Snapshot = Snapshot.(string)
	}
	configImportP := models.NewConfigurationImportPolicy(fmt.Sprintf("fabric/configimp-%s", name), "uni", desc, configImportPAttr)

	configImportP.Status = "modified"

	err := aciClient.Save(configImportP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_config_rs_import_source") {
		_, newRelParam := d.GetChange("relation_config_rs_import_source")
		err = aciClient.DeleteRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_config_rs_import_source")
		d.Partial(false)

	}
	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_trig_rs_triggerable")
		d.Partial(false)

	}
	if d.HasChange("relation_config_rs_remote_path") {
		_, newRelParam := d.GetChange("relation_config_rs_remote_path")
		err = aciClient.DeleteRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_config_rs_remote_path")
		d.Partial(false)

	}

	d.SetId(configImportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciConfigurationImportPolicyRead(d, m)

}

func resourceAciConfigurationImportPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	configImportP, err := getRemoteConfigurationImportPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setConfigurationImportPolicyAttributes(configImportP, d)

	configRsImportSourceData, err := aciClient.ReadRelationconfigRsImportSourceFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsImportSource %v", err)

	} else {
		d.Set("relation_config_rs_import_source", configRsImportSourceData)
	}

	trigRsTriggerableData, err := aciClient.ReadRelationtrigRsTriggerableFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation trigRsTriggerable %v", err)

	} else {
		d.Set("relation_trig_rs_triggerable", trigRsTriggerableData)
	}

	configRsRemotePathData, err := aciClient.ReadRelationconfigRsRemotePathFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsRemotePath %v", err)

	} else {
		d.Set("relation_config_rs_remote_path", configRsRemotePathData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciConfigurationImportPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "configImportP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
