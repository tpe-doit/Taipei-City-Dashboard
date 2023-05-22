<script setup>
import { onMounted, ref } from 'vue';
import { parseGuageData, sumAllData } from '../../assets/utilityFunctions/parseChartData'
const props = defineProps({
    content: Object, defaultType: {
        type: String,
        default: 'radialbar'
    }
})

const chartOptions = ref({
    chart: {
        plotBackgroundColor: null,
        plotBorderWidth: null,
        backgroundColor: null,
        plotShadow: false,
        inverted: true,
        polar: true,
        spaceTop: 0,
        height: '70%'
    },
    credits: {
        enabled: false,
    },
    exporting: {
        enabled: false,
    },
    legend: {
        enabled: false,
    },
    plotOptions: {
        column: {
            allowPointSelect: true,
            cursor: "pointer",
            dataLabels: {
                enabled: true,
                format: "{point.y}",
            },
            borderWidth: 2,
            borderColor: "#282a2c",
            pointWidth: 8,
            groupPadding: 4,
        },
        pie: {
            allowPointSelect: true,
            cursor: "pointer",
            dataLabels: {
                enabled: true,
                format: "<b>{point.name}</b>: {point.y}",
            },
            borderWidth: 2,
            borderColor: "#282a2c",
        },
    },
    title: {
        text: 'hello',
        verticalAlign: 'middle',
        style: {
            color: 'white',
            fontSize: '1rem'
        }
    },
    series: [
        {
            colorByPoint: true,
            innerSize: "75%",
            size: "50%",
            data: [],
        },
    ],
    xAxis: {
        text: null,
        type: 'category',
        labels: {
            formatter: function () {
                if (this.value.length > 6) {
                    return this.value.substring(0, 6) + '...';
                }
                return this.value;
            },
            style: {
                color: 'white',
            }
        }
    },
    yAxis: {
        text: null,
        type: 'category',
        visible: false,
    },
    tooltip: {
        pointFormat: "<b>{point.y}</b>",
        backgroundColor: "#090909",
        style: {
            color: "#888787",
        },
    },
})

onMounted(() => {
    chartOptions.value.series[0].data = parseGuageData(props.content.chartData[0].data)
    if (props.content.request_list[0].color && props.content.request_list[0].color.length > 1) {
        chartOptions.value.colors = props.content.request_list[0].color
    }

    if (props.defaultType === 'guage') {
        toGuage();
    } else if (props.defaultType === 'radialbar') {
        toRadialBar();
    }
})

function toGuage() {
    chartOptions.value.chart.type = 'pie'
}
function toRadialBar() {
    chartOptions.value.chart.type = 'column'
}

</script>

<template>
    <div class="guagedata">
        <highcharts :options="chartOptions" :style="{ height: '90%', marginTop: '1.5rem' }"></highcharts>
        <div class="guagedata-control">
            <button @click="toGuage">量表圖</button>
            <button @click="toRadialBar">環狀條形圖</button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.guagedata {
    min-height: 100%;
    max-height: 100%;
    overflow-y: scroll;

    &-control {
        position: absolute;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 0;

        button {
            background-color: rgb(77, 77, 77);
            padding: 4px 4px;
            border-radius: 5px;
            transition: color 0.2s, opacity 0.2s;
            font-size: var(--font-s);
            margin: 0 4px;
            color: var(--color-complement-text);
            opacity: 0.25;
            text-align: center;

            &:hover {
                color: white;
                opacity: 1;
            }
        }
    }

}
</style>