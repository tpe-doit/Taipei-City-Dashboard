<script setup>
import { onMounted, ref, watch } from 'vue';
import { parsePercentData, sumAllData } from '../../assets/utilityFunctions/parseChartData'
const props = defineProps({
    content: Object, defaultType: {
        type: String,
        default: 'pie'
    }
})

const chartOptions = ref({
    chart: {
        plotBackgroundColor: null,
        plotBorderWidth: null,
        backgroundColor: null,
        plotShadow: false,
        inverted: true,
        spaceTop: 0,
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
    chartOptions.value.series[0].data = parsePercentData(props.content.chartData[0].data)
    if (props.content.request_list[0].color && props.content.request_list[0].color.length > 1) {
        chartOptions.value.colors = props.content.request_list[0].color
    }

    if (props.defaultType === 'pie') {
        toPie();
    } else if (props.defaultType === 'column') {
        toColumn();
    }
})

function toPie() {
    chartOptions.value.chart.type = 'pie'
    chartOptions.value.chart.height = '70%'
    chartOptions.value.xAxis.visible = false;
    chartOptions.value.title.text = sumAllData(parsePercentData(props.content.chartData[0].data))

}
function toColumn() {
    chartOptions.value.chart.type = 'column'
    chartOptions.value.chart.height = props.content.chartData[0].data.length * 24
    chartOptions.value.xAxis.visible = true
    chartOptions.value.title.text = undefined
}

</script>

<template>
    <div class="percentdata">
        <highcharts :options="chartOptions" :style="{ height: '90%', marginTop: '1.5rem' }"></highcharts>
        <div class="percentdata-control">
            <button @click="toPie">圓餅圖</button>
            <button @click="toColumn">橫條圖</button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.percentdata {
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