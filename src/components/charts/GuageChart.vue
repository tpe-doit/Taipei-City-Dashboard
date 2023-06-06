<!-- Cleaned -->

<script setup>
import { ref, computed } from 'vue';

const props = defineProps(['chart_config', 'activeChart', 'series'])

// Guage charts in apexcharts uses a slightly different data format from other chart types
// As such, the following parsing function are required
const parseSeries = computed(() => {
    let output = {}
    let parsedSeries = []
    let parsedTooltip = []
    for (let i = 0; i < props.series[0].data.length; i++) {
        let total = props.series[0].data[i] + props.series[1].data[i]
        parsedSeries.push(Math.round(props.series[0].data[i] / total * 100))
        parsedTooltip.push(`${props.series[0].data[i]} / ${props.series[1].data[i]}`)
    }
    output.series = parsedSeries
    output.tooltipText = parsedTooltip
    return output
})

// chartOptions needs to be in the bottom since it uses computed data
const chartOptions = ref({
    chart: {
        toolbar: {
            show: false,
        },
    },
    colors: props.chart_config.color,
    labels: props.chart_config.categories ? props.chart_config.categories : [],
    legend: {
        show: false,
    },
    plotOptions: {
        radialBar: {
            dataLabels: {
                name: {
                    color: '#888787',
                    fontSize: '0.8rem',
                },
                total: {
                    color: '#888787',
                    fontSize: '0.8rem',
                    label: '平均',
                    show: true,
                },
                value: {
                    color: '#888787',
                    fontSize: '16px',
                    offsetY: 5,
                },
            },
            track: {
                background: "#777"
            },
        }
    },
    tooltip: {
        custom: function ({ seriesIndex, w }) {
            // The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.seriesNames[seriesIndex] + '</h6>' +
                '<span>' + `${parseSeries.value.tooltipText[seriesIndex]}` + '</span>' +
                '</div>'
        },
        enabled: true,
    },
})
</script>

<template>
    <div v-if="activeChart === 'GuageChart'">
        <apexchart width="80%" height="300px" type="radialBar" :options="chartOptions" :series="parseSeries.series">
        </apexchart>
    </div>
</template>