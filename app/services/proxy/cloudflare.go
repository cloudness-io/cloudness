package proxy

import (
	"context"

	"github.com/cloudness-io/cloudness/app/usererror"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

type cloudflareProxy struct {
}

func newCloudflareProxy() *cloudflareProxy {
	return &cloudflareProxy{}
}

func (p *cloudflareProxy) ValidateAPIKeyForDNS01(ctx context.Context, apitoken string, zone string) error {
	api, err := cloudflare.NewWithAPIToken(apitoken)
	if err != nil {
		return err
	}

	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any("zone", zone).Msg("Proxy: get zone id call failed")
		return usererror.Forbidden("Zone:Read permission missing or zone not found")
	}

	dnsRecord, err := api.CreateDNSRecord(ctx, &cloudflare.ResourceContainer{Identifier: zoneID},
		cloudflare.CreateDNSRecordParams{Type: "TXT", Name: "proxy." + zone, Content: "This is a test record. If you see this, please feel free to delete it."},
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Proxy: reading dns records failed")
		return usererror.Forbidden("Zone:Write permission missing")
	}

	if err := api.DeleteDNSRecord(ctx, &cloudflare.ResourceContainer{Identifier: zoneID}, dnsRecord.ID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Proxy: error deleting dns record")
	}

	return nil
}
