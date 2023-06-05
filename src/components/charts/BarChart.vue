<script setup>
import { ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const options = ref({
    legend: {
        show: false
    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 0,
    },
    tooltip: {
        followCursor: true,
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[dataPointIndex] + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        textAnchor: 'start',
        offsetX: 20
    },
    plotOptions: {
        bar: {
            horizontal: true,
            borderRadius: 2,
            distributed: true
        }
    },
    colors: props.chart_config.color,
    grid: {
        show: false,
    },
    chart: {
        toolbar: {
            show: false
        },
        stacked: true,
        offsetY: 15,
    },
    xaxis: {
        type: 'category',
        labels: {
            show: false
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false
        }
    },
    yaxis: {
        labels: {
            formatter: function (value) {
                return value.length > 7 ? value.slice(0, 6) + "..." : value
            }
        }
    }
})

const height = computed(() => {
    return `${40 + props.series[0].data.length * 16}`
})

</script>

<template>
    <div v-if="activeChart === 'BarChart'">
        <apexchart width="100%" :height="height" type="bar" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped></style>