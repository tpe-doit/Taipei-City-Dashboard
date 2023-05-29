<script setup>
import { computed, ref } from 'vue';
import DialogContainer from './DialogContainer.vue';
import ComponentTag from '../utilities/ComponentTag.vue'
import { useDialogStore } from '../../store/dialogStore';
import { useContentStore } from '../../store/contentStore';

import CustomCheckBox from '../utilities/CustomCheckBox.vue';

const dialogStore = useDialogStore()
const contentStore = useContentStore()

const outputList = computed(() => {

    let output = [...availableComponents.value]

    if (searchName.value) {
        output = output.filter((item) => item.name.includes(searchName.value))
    }
    if (searchIndex.value) {
        output = output.filter((item) => item.index.toString().includes(searchIndex.value))
    }
    if (filterSource.value.length > 0) {
        output = output.filter((item) => filterSource.value.includes(item.source))
    }
    if (filterControl.value.includes('有地理資料')) {
        output = output.filter((item) => item.map_config !== null)
    }

    return output
})
const availableComponents = computed(() => {
    const allComponentIds = Object.keys(contentStore.components)
    const taken = contentStore.currentDashboard.content.map((item) => item.id)
    const maplayer = contentStore.mapLayers.map((item) => item.id)
    const available = allComponentIds.filter(item => !taken.includes(+item) && !maplayer.includes(+item))
    const output = []
    available.forEach(element => {
        output.push(contentStore.components[element])
    });
    return output
})

const searchName = ref('');
const searchIndex = ref('');
const filterSource = ref([])
const filterType = ref([])
const filterFrequency = ref([])
const filterControl = ref([])

const componentsSelected = ref([])

const filterItems = {
    source: ['交通局', '警察局', '都發局', '消防局', '社會局', '工務局'],
    type: ['交通', '產業', '土地', '安全'],
    frequency: ['無定期更新', '每半年', '每個月', '每兩週', '每一週', '每一天', '每一小時'],
    control: ['有地理資料']
}

function dataTime(time_from, time_to) {
    if (!time_from) {
        return '固定資料'
    }
    if (!time_to) {
        return time_from.slice(0, 10)
    }
    return `${time_from.slice(0, 10)} ~ ${time_to.slice(0, 10)}`
}

function updateFreq(update_freq, update_freq_unit) {
    const unitRef = {
        minute: "分鐘",
        hour: "小時",
        day: "天",
        week: "週",
        month: "月",
        year: "年"
    }
    if (!update_freq) {
        return '不定期更新'
    }
    return `每${update_freq}${unitRef[update_freq_unit]}更新`
}

// Uncomment the following to restore add components function if new backend is connected
function handleSubmit() {
    // contentStore.addComponents(componentsSelected.value)
    dialogStore.showNotification('fail', '尚未新增新增儀表板功能，無法新增組件')
    handleClose();
}
function handleClose() {
    searchName.value = ''
    searchIndex.value = ''
    componentsSelected.value = []
    clearFilters()
    dialogStore.hideAllDialogs()
}
function clearFilters() {
    filterSource.value = []
    filterControl.value = []
    filterType.value = []
    filterFrequency.value = []
}

</script>

<template>
    <DialogContainer :dialog="`addComponent`" @onClose="handleClose">
        <div class="addcomponent">
            <div class="addcomponent-header">
                <h2>新增組件</h2>
                <div class="addcomponent-header-search">
                    <div class="flex">
                        <div :style="{ position: 'relative' }">
                            <input type="text" placeholder="以組件名稱搜尋" v-model="searchName" />
                            <span class="clear" v-if="searchName" @click="() => { searchName = '' }">cancel</span>
                        </div>
                        <div :style="{ position: 'relative' }">
                            <input type="text" placeholder="以組件Index搜尋" v-model="searchIndex" />
                            <span class="clear" v-if="searchIndex" @click="() => { searchIndex = '' }">cancel</span>
                        </div>
                    </div>
                    <div class="flex">
                        <button @click="handleClose">取消</button>
                        <button v-if="componentsSelected.length > 0"
                            @click="handleSubmit"><span>add_chart</span>確認新增</button>
                    </div>
                </div>
            </div>
            <div class="addcomponent-filter">
                <button @click="clearFilters">清除篩選條件</button>
                <div>
                    <h3>依資料來源篩選</h3>
                    <div v-for="item in filterItems.source" :key="item">
                        <input type="checkbox" :id="item" :value="item" v-model="filterSource" class="custom-check-input" />
                        <CustomCheckBox :for="item">{{ item }}</CustomCheckBox>
                    </div>
                    <h3>依類別標籤篩選</h3>
                    <div v-for="item in filterItems.type" :key="item">
                        <input type="checkbox" :id="item" :value="item" v-model="filterType" class="custom-check-input" />
                        <CustomCheckBox :for="item">{{ item }}</CustomCheckBox>
                    </div>
                    <h3>依更新頻率篩選</h3>
                    <div v-for="item in filterItems.frequency" :key="item">
                        <input type="checkbox" :id="item" :value="item" v-model="filterFrequency"
                            class="custom-check-input" />
                        <CustomCheckBox :for="item">{{ item }}</CustomCheckBox>
                    </div>
                    <h3>依功能種類篩選</h3>
                    <div>
                        <div v-for="item in filterItems.control" :key="item">
                            <input type="checkbox" :id="item" :value="item" v-model="filterControl"
                                class="custom-check-input" />
                            <CustomCheckBox :for="item">{{ item }}</CustomCheckBox>
                        </div>
                    </div>
                </div>
            </div>
            <div>
                <p :style="{ margin: '0 0 0.5rem 1rem' }">計 {{ outputList.length }} 個組件符合篩選條件 | 共選取 {{
                    componentsSelected.length }} 個</p>
                <div class="addcomponent-list">
                    <div v-for="item in outputList" :key="item.id">
                        <input type="checkbox" :id="item.name" :value="item.id" v-model="componentsSelected"
                            :style="{ display: 'none' }" />
                        <label :for="item.name" class="addcomponent-list-item">
                            <div>
                                我是縮圖
                            </div>
                            <div>
                                <div>
                                    <h2>{{ item.name }}</h2>
                                    <h3>{{ `${item.source} | ${dataTime(item.time_from, item.time_to)}` }}</h3>
                                </div>
                                <div :style="{ margin: '0.75rem 0' }">
                                    <p>{{ item.short_desc }}</p>
                                    <br />
                                    <p><b>組件ID：</b>{{ item.id }}</p>
                                    <p><b>組件Index：</b>{{ item.index }}</p>
                                </div>
                                <div :style="{ display: 'flex', marginTop: '0.5rem' }">
                                    <ComponentTag v-for="element in item.tags" icon="" :text="element" mode="fill" />
                                </div>
                                <div :style="{ display: 'flex', marginTop: '4px' }">
                                    <ComponentTag :icon="``"
                                        :text="`${updateFreq(item.update_freq, item.update_freq_unit)}`" />
                                    <ComponentTag v-if="item.map_config === null" :icon="`wrong_location`" :text="`沒有地圖`" />
                                </div>
                            </div>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    </DialogContainer>
</template>

<style scoped lang="scss">
.addcomponent {
    display: grid;
    width: 800px;
    height: 600px;
    grid-template-columns: 180px 1fr;
    grid-template-areas:
        "header header"
        "filter list";
    grid-template-rows: 5rem 1fr;

    &-header {
        grid-area: header;

        h2 {
            font-size: var(--font-l);
        }

        &-search {
            display: flex;
            justify-content: space-between;
            margin-top: 0.5rem;

            input {
                margin-right: 0.5rem;
                width: 200px;
            }

            span {
                font-family: var(--font-icon);
                font-size: calc(var(--font-m) * var(--font-to-icon));
                margin-right: 4px
            }

            button {
                display: flex;
                align-items: center;
                margin-right: 0.4rem;
                font-size: var(--font-m);
                border-radius: 5px;
                justify-self: baseline;

                &:nth-child(2) {
                    background-color: var(--color-highlight);
                    padding: 2px 4px;
                }
            }

            .clear {
                font-size: var(--font-m);
                color: var(--color-complement-text);
                position: absolute;
                right: 0.5rem;
                top: 0.4rem;
                cursor: pointer;
                transition: color 0.2s;

                &:hover {
                    color: var(--color-highlight)
                }
            }
        }
    }

    &-filter {
        grid-area: filter;
        border-right: solid 1px var(--color-border);
        padding-right: 1rem;
        overflow-y: scroll;

        span {
            font-family: var(--font-icon);
            font-size: calc(var(--font-m) * var(--font-to-icon));
        }

        button {
            color: var(--color-highlight);
        }

        h3 {
            font-size: var(--font-m);
            font-weight: 400;
            color: var(--color-complement-text);
            margin: 0.5rem 0;
        }

        label {
            margin-left: 0.5rem;
            font-size: var(--font-m);
        }

        input {
            display: none;
        }
    }

    &-list {
        grid-area: list;
        overflow-y: scroll;
        padding-left: 1rem;
        max-height: calc(100% - 1rem);

        &-item {
            width: calc(100% - 1.2rem);
            border-radius: 5px;
            border: solid 1px var(--color-border);
            padding: 0.5rem;
            margin: 0.5rem 0;
            display: grid;
            cursor: pointer;
            grid-template-columns: 150px 1fr;
            column-gap: 1rem;

            transition: border-color 0.2s,
                border-width 0.2s;

            h3 {
                font-weight: 400;
                color: var(--color-complement-text)
            }

            >div:first-child {
                background-color: rgb(75, 75, 75);
                border-radius: 5px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
        }

        input:checked+&-item {
            border-color: var(--color-highlight);
        }
    }
}

.flex {
    display: flex;
    justify-content: space-between;
}
</style>