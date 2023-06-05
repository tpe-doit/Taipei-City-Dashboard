<script setup>
import { ref } from 'vue';
import { useMapStore } from '../../store/mapStore';

const mapStore = useMapStore()

const props = defineProps(['content'])
const checked = ref(false)

function handleToggle() {
    if (!props.content.map_config) {
        return
    }
    if (checked.value) {
        mapStore.addToMapLayerList(props.content.map_config)
    } else {
        mapStore.turnOffMapLayerVisibility(props.content.map_config)
    }
}

</script>

<template>
    <div class="mobilelayertab">
        <input :id="content.index" type="checkbox" v-model="checked" @change="handleToggle" />
        <label :for="content.index"></label>
        <p>{{ content.name.length > 6 ? `${content.name.slice(0, 5)}...` : content.name }}</p>
    </div>
</template>

<style scoped lang="scss">
.mobilelayertab {
    input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    label {
        display: inline-block;
        width: 73px;
        height: 73px;
        border: solid 1px transparent;
        border-radius: 5px;
        cursor: pointer;
        background-color: rgb(78, 71, 71);
        transition: border 0.2s;
    }

    input:checked+label {
        border: solid 1px var(--color-highlight);
    }

    p {
        color: var(--color-complement-text);
        font-size: 0.75rem;
        text-align: center;
    }

    margin-bottom: 8px
}
</style>