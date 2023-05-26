<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

function parseTime(time) {
    return time.replace('T', ' ').slice(0, -4)
}

const options = ref({
    legend: {
        show: props.series.length > 1 ? true : false
    },
    stroke: {
        show: true,
        colors: props.chart_config.color,
        width: 2,
        curve: 'smooth'
    },
    markers: {
        size: 3,
        strokeWidth: 0,
        hover: {
            size: 5,
        }
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + `${parseTime(w.config.series[seriesIndex].data[dataPointIndex].x)}` + ` - ${w.globals.seriesNames[seriesIndex]}` + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        enabled: false
    },
    colors: props.chart_config.color,
    grid: {
        show: false,
    },
    chart: {
        toolbar: {
            show: false,
            tools: {
                zoom: false,
            }
        },
        stacked: true
    },
    xaxis: {
        type: 'datetime',
        axisTicks: {
            show: false,
        },
        crosshairs: {
            show: false
        },
        tooltip: {
            enabled: false
        }
    },
})

</script>

<template>
    <div v-if="activeChart === 'TimelineStackedChart'">
        <apexchart width="100%" height="260px" type="area" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped></style>