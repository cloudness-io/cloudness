// Alpine.js component for metrics chart
function metricsChart() {
    return {
        config: {},
        chart: null,

        initChart() {
            const container = this.$el;

            // Prepare data in uPlot format: [timestamps, ...series values]
            const seriesData = this.config.series || [];
            // Use timestamps from first series (assuming all series share same timestamps)
            const utcTimestamps = seriesData.length > 0 ? (seriesData[0].timestamps || []) : [];

            // Convert UTC timestamps to local timezone
            const timestamps = utcTimestamps.map(ts => {
                // If timestamp is a string, parse it as UTC and convert to local Unix timestamp
                if (typeof ts === 'string') {
                    return new Date(ts + (ts.endsWith('Z') ? '' : 'Z')).getTime() / 1000;
                }
                // If already a Unix timestamp, return as-is (uPlot will handle timezone display)
                return ts;
            });

            const data = [
                timestamps,
                ...seriesData.map(s => (s.values || []).map(v => v === 0 ? null : v))
            ];

            // Color palette for different series
            const colors = ['#217ecaff', '#0b7954ff', '#a36c0cff', '#812828ff', '#4f3392ff', '#a31d60ff'];

            // Build series config
            const seriesConfig = [
                { label: 'Time' },
                ...seriesData.map((s, i) => ({
                    label: s.label || `Series ${i + 1}`,
                    stroke: colors[i % colors.length],
                    width: 2,
                    points: { show: false }
                }))
            ];

            const opts = {
                title: this.config.title || 'Metrics',
                width: container.offsetWidth,
                height: this.config.height || 300,
                series: seriesConfig,
                axes: [
                    {
                        space: 80,
                        stroke: '#989694',
                        grid: {
                            show: false,
                        },
                    },
                    {
                        space: 40,
                        stroke: '#989694',
                        grid: {
                            show: false,
                        },
                        values: (u, vals) => vals.map(v =>
                            v.toFixed(this.config.valuePrecision || 2)
                        ),
                    }
                ],
                scales: {
                    x: {
                        time: true,
                    },
                },
                legend: {
                    show: true,
                    live: true
                },
                cursor: {
                    drag: {
                        x: false,
                        y: false
                    },
                    focus: {
                        prox: 30
                    },
                    points: {
                        size: 8,
                        width: 2,
                        stroke: (u, i) => colors[(i - 1) % colors.length],
                        fill: '#242423'
                    }
                },
            };

            this.chart = new uPlot(opts, data, container);

            // Handle resize
            window.addEventListener('resize', () => this.handleResize());

            // Clean up on destroy
            this.$el.addEventListener('alpine:destroying', () => this.destroy());
        },

        handleResize() {
            if (this.chart) {
                this.chart.setSize({
                    width: this.$el.offsetWidth,
                    height: this.config.height || 300
                });
            }
        },

        updateData(newTimestamps, newValues) {
            if (this.chart) {
                this.chart.setData([newTimestamps, newValues]);
            }
        },

        destroy() {
            if (this.chart) {
                this.chart.destroy();
                this.chart = null;
            }
        }
    };
}