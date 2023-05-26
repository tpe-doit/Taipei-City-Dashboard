<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

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

const options = ref({
    legend: {
        show: false
    },
    tooltip: {
        enabled: true,
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.seriesNames[seriesIndex] + '</h6>' +
                '<span>' + `${parseSeries.value.tooltipText[seriesIndex]}` + '</span>' +
                '</div>'
        }
    },
    plotOptions: {
        radialBar: {
            track: {
                background: "#777"
            },
            dataLabels: {
                name: {
                    fontSize: '0.8rem',
                    color: '#888787',
                },
                value: {
                    fontSize: '16px',
                    color: '#888787',
                    offsetY: 5,
                },
                total: {
                    show: true,
                    fontSize: '0.8rem',
                    color: '#888787',
                    label: '平均'
                },
            }

        }
    },
    colors: props.chart_config.color,
    chart: {
        toolbar: {
            show: false
        },
    },
    labels: props.chart_config.categories ? props.chart_config.categories : []
})

</script>

<template>
    <div v-if="activeChart === 'GuageChart'" class="guagechart">
        <apexchart width="80%" height="300px" type="radialBar" :options="options" :series="parseSeries.series">
        </apexchart>
    </div>
</template>

<style scoped lang="scss"></style>