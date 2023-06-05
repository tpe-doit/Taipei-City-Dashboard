<script setup>
import { ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const options = ref({
    legend: {
        show: props.chart_config.categories ? true : false,

    },
    stroke: {
        show: true,
        width: 2,
    },
    markers: {
        size: 3,
        strokeWidth: 0,
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[dataPointIndex] + `${props.chart_config.categories ? '-' + w.globals.seriesNames[seriesIndex] : ''}` + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        }
    },
    plotOptions: {
        radar: {
            polygons: {
                strokeColors: '#555',
                connectorColors: '#444'
            }
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
    },
    xaxis: {
        type: 'category',
        categories: props.chart_config.categories ? props.chart_config.categories : [],
        labels: {
            offsetY: 5,
            formatter: function (value) {
                return value.length > 7 ? value.slice(0, 6) + "..." : value
            }
        }
    },
    yaxis: {
        labels: {
            formatter: (value) => { return '' },
        },
        axisBorder: {
            color: '#000',
        }
    }
})

</script>

<template>
    <div v-if="activeChart === 'RadarChart'">
        <apexchart width="100%" height="270px" type="radar" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped></style>