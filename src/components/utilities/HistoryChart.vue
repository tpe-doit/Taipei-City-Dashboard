<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'series'])

const options = ref({
    legend: {
        show: false
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
                '<h6>' + `${w.config.series[seriesIndex].data[dataPointIndex].x}` + '</h6>' +
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
            tools: {
                zoomin: false,
                zoomout: false,
                pan: false,
                download: false,
                reset: '<p>' + '重置' + '</p>'
            }
        }
    },
    xaxis: {
        type: 'datetime',
        crosshairs: {
            show: false
        },
        tooltip: {
            enabled: false
        },
    },
})

</script>

<template>
    <div>
        <apexchart width="100%" height="140px" type="area" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped lang="scss"></style>