<script setup>
import { defineProps, ref, computed } from 'vue';
import MetroCarDensity from '../utilities/MetroCarDensity.vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const lineInfo = {
    "BL": [
        { "name": "南港展覽館", "id": "23" },
        { "name": "南港", "id": "22" },
        { "name": "昆陽", "id": "21" },
        { "name": "後山埤", "id": "20" },
        { "name": "永春", "id": "19" },
        { "name": "市政府", "id": "18" },
        { "name": "國父紀念館", "id": "17" },
        { "name": "忠孝敦化", "id": "16" },
        { "name": "忠孝復興", "id": "15" },
        { "name": "忠孝新生", "id": "14" },
        { "name": "善導寺", "id": "13" },
        { "name": "台北車站", "id": "12" },
        { "name": "西門", "id": "11" },
        { "name": "龍山寺", "id": "10" },
        { "name": "江子翠", "id": "09" },
        { "name": "新埔", "id": "08" },
        { "name": "板橋", "id": "07" },
        { "name": "府中", "id": "06" },
        { "name": "亞東醫院", "id": "05" },
        { "name": "海山", "id": "04" },
        { "name": "土城", "id": "03" },
        { "name": "永寧", "id": "02" },
        { "name": "頂埔", "id": "01" }
    ],
    "R": [
        { "name": "淡水", "id": "28" },
        { "name": "紅樹林", "id": "27" },
        { "name": "竹圍", "id": "26" },
        { "name": "關渡", "id": "25" },
        { "name": "忠義", "id": "24" },
        { "name": "復興崗", "id": "23" },
        { "name": "北投", "id": "22" },
        { "name": "奇岩", "id": "21" },
        { "name": "唭哩岸", "id": "20" },
        { "name": "石牌", "id": "19" },
        { "name": "明德", "id": "18" },
        { "name": "芝山", "id": "17" },
        { "name": "士林", "id": "16" },
        { "name": "劍潭", "id": "15" },
        { "name": "圓山", "id": "14" },
        { "name": "民權西路", "id": "13" },
        { "name": "雙連", "id": "12" },
        { "name": "中山", "id": "11" },
        { "name": "台北車站", "id": "10" },
        { "name": "台大醫院", "id": "09" },
        { "name": "中正紀念堂", "id": "08" },
        { "name": "東門", "id": "07" },
        { "name": "大安森林公園", "id": "06" },
        { "name": "大安", "id": "05" },
        { "name": "信義安和", "id": "04" },
        { "name": "台北101/世貿", "id": "03" },
        { "name": "象山", "id": "02" }
    ]
}
const color = props.chart_config.color[0]
const line = props.series[0].name
</script>

<template>
    <div v-if="activeChart === 'MetroChart'" class="metrochart">
        <div v-for="(item, index) in  lineInfo[line] " :class="`initial-animation-${index + 1}`">
            <div class="metrochart-block">
                <h5>{{ item.name }}</h5>
                <div class="metrochart-block-tag"
                    :style="{ borderColor: color, backgroundColor: (index === lineInfo[line].length - 1 || index === 0) ? color : 'white' }">
                    <p :style="{ color: (index === lineInfo[line].length - 1 || index === 0) ? 'white' : 'black' }">{{ line
                    }}</p>
                    <p :style="{ color: (index === lineInfo[line].length - 1 || index === 0) ? 'white' : 'black' }">{{
                        item.id }}</p>
                </div>
                <MetroCarDensity :weight="series[1].data.find(element => element.x === item.id)" direction="south" />
                <MetroCarDensity :weight="series[0].data.find(element => element.x === item.id)" direction="north" />
            </div>
            <div class="metrochart-block">
                <div></div>
                <div class="metrochart-block-line"
                    :style="{ backgroundColor: index === lineInfo[line].length - 1 ? 'transparent' : color }"></div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.metrochart {
    width: 100%;

    h5 {
        font-size: 0.7rem;
        margin-right: 0.4rem;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        pointer-events: none;
        user-select: none;
        font-weight: 400;
    }

    p {
        font-size: 0.6rem;
        color: black;
        pointer-events: none;
        user-select: none;
    }

    &-block {
        display: grid;
        grid-template-columns: 5rem 20px 1fr 1fr;

        &-tag {
            width: 1rem;
            background-color: white;
            height: 1.4rem;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            border-width: 2px;
            border-style: solid;
            border-radius: 4px;
        }

        &-line {
            margin: 0px 6px;
            width: 8px;
            height: 1rem;
        }
    }
}

@keyframes ease-in {
    0% {
        opacity: 0
    }

    ;

    100% {
        opacity: 1
    }
}

@for $i from 1 through 40 {
    .initial-animation-#{$i} {
        animation-name: ease-in;
        animation-duration: 0.2s;
        animation-delay: 0.05s * ($i - 1);
        animation-timing-function: linear;
        animation-fill-mode: forwards;
        opacity: 0;
    }
}
</style>