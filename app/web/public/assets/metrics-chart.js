// Alpine.js component for metrics chart using TradingView Lightweight Charts
// Each series is rendered as an individual area series, trimmed to its active
// range so terminated instances don't trail zeros and new instances don't show
// leading zeros before deployment.

function metricsChart() {
    return {
        config: {},
        chart: null,
        seriesInstances: [],
        _resizeHandler: null,

        initChart() {
            const container = this.$el;
            const seriesData = this.config.series || [];

            // Color palette: line + translucent area fill
            const colors = [
                { line: '#217eca', area: 'rgba(33, 126, 202, 0.25)' },
                { line: '#0b7954', area: 'rgba(11, 121, 84, 0.25)' },
                { line: '#a36c0c', area: 'rgba(163, 108, 12, 0.25)' },
                { line: '#812828', area: 'rgba(129, 40, 40, 0.25)' },
                { line: '#4f3392', area: 'rgba(79, 51, 146, 0.25)' },
                { line: '#a31d60', area: 'rgba(163, 29, 96, 0.25)' },
            ];

            // Add title above chart
            if (this.config.title) {
                const titleEl = document.createElement('div');
                titleEl.textContent = this.config.title;
                titleEl.style.cssText = 'color: #989694; font-size: 14px; margin-bottom: 8px; font-weight: 500;';
                container.appendChild(titleEl);
            }

            // Create a sub-container for the chart canvas
            const chartEl = document.createElement('div');
            container.appendChild(chartEl);
            this._chartEl = chartEl;

            const chartHeight = (this.config.height || 300) - (this.config.title ? 30 : 0);

            this.chart = LightweightCharts.createChart(chartEl, {
                width: container.offsetWidth,
                height: chartHeight,
                layout: {
                    background: { type: 'solid', color: 'transparent' },
                    textColor: '#989694',
                    fontSize: 12,
                    attributionLogo: false,
                },
                grid: {
                    vertLines: { visible: false },
                    horzLines: { visible: false },
                },
                crosshair: {
                    mode: LightweightCharts.CrosshairMode.Normal,
                },
                rightPriceScale: {
                    borderVisible: false,
                    scaleMargins: { top: 0.05, bottom: 0.05 },
                },
                timeScale: {
                    borderVisible: false,
                    timeVisible: true,
                    secondsVisible: false,
                },
                localization: {
                    priceFormatter: (price) => price.toFixed(this.config.valuePrecision || 2),
                },
            });

            // Add an area series for each data series, trimmed to its active range
            seriesData.forEach((s, i) => {
                const data = this._buildTrimmedSeriesData(s);
                if (data.length === 0) return; // skip entirely empty series

                const color = colors[i % colors.length];
                const areaSeries = this.chart.addAreaSeries({
                    lineColor: color.line,
                    topColor: color.area,
                    bottomColor: 'transparent',
                    lineWidth: 2,
                    title: s.label || `Series ${i + 1}`,
                    crosshairMarkerRadius: 4,
                    crosshairMarkerBorderWidth: 2,
                    crosshairMarkerBorderColor: color.line,
                    crosshairMarkerBackgroundColor: '#242423',
                    lastValueVisible: false,
                    priceLineVisible: false,
                });

                areaSeries.setData(data);
                this.seriesInstances.push(areaSeries);
            });

            this.chart.timeScale().fitContent();

            // Handle resize
            this._resizeHandler = () => this.handleResize();
            window.addEventListener('resize', this._resizeHandler);

            // Clean up on destroy
            this.$el.addEventListener('alpine:destroying', () => this.destroy());
        },

        // Trim leading and trailing zeros from a series so that:
        // - Terminated instances stop at their last reported value
        // - New instances start from their first reported value
        // Zeros in the middle (brief dips) are preserved.
        _buildTrimmedSeriesData(series) {
            const timestamps = series.timestamps || [];
            const values = series.values || [];

            // Find the active range (first and last non-zero value)
            let firstNonZero = -1;
            let lastNonZero = -1;
            for (let i = 0; i < values.length; i++) {
                if (values[i] !== 0) {
                    if (firstNonZero === -1) firstNonZero = i;
                    lastNonZero = i;
                }
            }

            if (firstNonZero === -1) return []; // entire series is zeros

            // Offset to shift UTC timestamps to local time for display
            const tzOffsetSec = new Date().getTimezoneOffset() * -60;

            const data = [];
            for (let j = firstNonZero; j <= lastNonZero; j++) {
                const ts = timestamps[j];
                const utc = typeof ts === 'string'
                    ? Math.floor(new Date(ts + (ts.endsWith('Z') ? '' : 'Z')).getTime() / 1000)
                    : ts;

                data.push({ time: utc + tzOffsetSec, value: values[j] });
            }

            return data;
        },

        handleResize() {
            if (this.chart) {
                const chartHeight = (this.config.height || 300) - (this.config.title ? 30 : 0);
                this.chart.applyOptions({
                    width: this.$el.offsetWidth,
                    height: chartHeight,
                });
            }
        },

        updateData(newSeriesData) {
            // Remove old series and re-add with new data
            if (this.chart) {
                this.seriesInstances.forEach(s => this.chart.removeSeries(s));
                this.seriesInstances = [];

                const colors = [
                    { line: '#217eca', area: 'rgba(33, 126, 202, 0.25)' },
                    { line: '#0b7954', area: 'rgba(11, 121, 84, 0.25)' },
                    { line: '#a36c0c', area: 'rgba(163, 108, 12, 0.25)' },
                    { line: '#812828', area: 'rgba(129, 40, 40, 0.25)' },
                    { line: '#4f3392', area: 'rgba(79, 51, 146, 0.25)' },
                    { line: '#a31d60', area: 'rgba(163, 29, 96, 0.25)' },
                ];

                newSeriesData.forEach((s, i) => {
                    const data = this._buildTrimmedSeriesData(s);
                    if (data.length === 0) return;

                    const color = colors[i % colors.length];
                    const areaSeries = this.chart.addAreaSeries({
                        lineColor: color.line,
                        topColor: color.area,
                        bottomColor: 'transparent',
                        lineWidth: 2,
                        title: s.label || `Series ${i + 1}`,
                        lastValueVisible: false,
                        priceLineVisible: false,
                    });

                    areaSeries.setData(data);
                    this.seriesInstances.push(areaSeries);
                });

                this.chart.timeScale().fitContent();
            }
        },

        destroy() {
            if (this._resizeHandler) {
                window.removeEventListener('resize', this._resizeHandler);
            }
            if (this.chart) {
                this.chart.remove();
                this.chart = null;
                this.seriesInstances = [];
            }
        }
    };
}