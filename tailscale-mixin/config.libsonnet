{
  _config+:: {
    local this = self,

    // Default datasource name
    datasourceName: 'default',

    // Opt-in to multiCluster dashboards by overriding this and the clusterLabel.
    showMultiCluster: false,
    clusterLabel: 'cluster',

    tailscaleSelector: 'job="tailscale-exporter"',
    // This selector is anything for now as scraping machines can vary in label names.
    tailscaledSelector: 'job=~".*"',

    // Tailnet
    tailscaleDeviceUnauthorizedEnabled: true,
    tailscaleDeviceUnauthorizedFor: '15m',
    tailscaleDeviceUnauthorizedSeverity: 'warning',

    tailscaleUserUnapprovedEnabled: true,
    tailscaleUserUnapprovedFor: '15m',
    tailscaleUserUnapprovedSeverity: 'warning',

    tailscaleUserRecentlyCreatedEnabled: true,
    tailscaleUserRecentlyCreatedSeverity: 'info',
    tailscaleUserRecentlyCreatedThreshold: '300',  // Seconds

    tailscaleDeviceUnapprovedRoutesEnabled: true,
    tailscaleDeviceUnapprovedRoutesFor: '15m',
    tailscaleDeviceUnapprovedRoutesSeverity: 'warning',
    tailscaleDeviceUnapprovedRoutesThreshold: '10',

    // Tailscaled
    tailscaledMachineHighOutboundDroppedPacketsEnabled: true,
    tailscaledMachineHighOutboundDroppedPacketsFor: '15m',
    tailscaledMachineHighOutboundDroppedPacketsSeverity: 'warning',
    tailscaledMachineHighOutboundDroppedPacketsThreshold: '50',

    grafanaUrl: 'https://grafana.com',

    dashboardIds: {
      'tailscale-overview': 'tailscale-mixin-over-k12e',
      'tailscale-machine': 'tailscaled-mixin-over-k12e',
    },
    dashboardUrls: {
      'tailscale-overview': '%s/d/%s/tailscale-overview' % [this.grafanaUrl, this.dashboardIds['tailscale-overview']],
      'tailscale-machine': '%s/d/%s/tailscale-machine' % [this.grafanaUrl, this.dashboardIds['tailscale-machine']],
    },

    tags: ['tailscale', 'tailscale-mixin'],

    // Custom annotations to display in graphs
    annotation: {
      enabled: false,
      name: 'Custom Annotation',
      tags: [],
      datasource: '-- Grafana --',
      iconColor: 'blue',
      type: 'tags',
    },
  },
}
