package migrations

import "github.com/goravel/framework/facades"

type M20251118152211UpdateAdsLogDetailsTable struct{}

func (r *M20251118152211UpdateAdsLogDetailsTable) Signature() string {
	return "20251118152211_update_ads_log_details_table"
}

func (r *M20251118152211UpdateAdsLogDetailsTable) Up() error {
	sql := `
ALTER TABLE ads_log_details
  ADD INDEX idx_ald_campaign_created_at (ads_campaign_id, created_at),
  ADD INDEX idx_ald_campaign_country (ads_campaign_id, country),
  ADD INDEX idx_ald_campaign_city (ads_campaign_id, city),
  ADD INDEX idx_ald_campaign_region (ads_campaign_id, region),
  ADD INDEX idx_ald_campaign_os (ads_campaign_id, os_name, os_version),
  ADD INDEX idx_ald_campaign_device (ads_campaign_id, device_type, device_name),
  ADD INDEX idx_ald_campaign_browser (ads_campaign_id, browser_name),
  ADD INDEX idx_ald_campaign_ip_useragent (ads_campaign_id, ip, user_agent),
  ADD FULLTEXT INDEX ft_ald_search (
    ip, country, region, city, user_agent, referrer, device_type, device_name, os_name, browser_name
  ),
  ALGORITHM=INPLACE,
  LOCK=SHARED;
`
	return facades.DB().Statement(sql)
}

func (r *M20251118152211UpdateAdsLogDetailsTable) Down() error {
	sql := `
ALTER TABLE ads_log_details
  DROP INDEX idx_ald_campaign_created_at,
  DROP INDEX idx_ald_campaign_country,
  DROP INDEX idx_ald_campaign_city,
  DROP INDEX idx_ald_campaign_region,
  DROP INDEX idx_ald_campaign_os,
  DROP INDEX idx_ald_campaign_device,
  DROP INDEX idx_ald_campaign_browser,
  DROP INDEX idx_ald_campaign_ip_useragent,
  DROP INDEX ft_ald_search;
`
	return facades.DB().Statement(sql)
}
