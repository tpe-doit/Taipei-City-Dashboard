<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useContentStore } from '../../store/contentStore';

const route = useRoute()
const contentStore = useContentStore()
const props = defineProps({
    icon: { type: String },
    title: { type: String },
    index: { type: String },
    expanded: { type: Boolean }
})

const tabLink = computed(() => {
    const path = route.path
    return `${route.path}?index=${props.index}`
})
const linkActiveOrNot = computed(() => {
    return route.query.index === props.index ? true : false;
})

</script>

<template>
    <router-link :to="tabLink" :class="{ sidebartab: true, 'link-active': linkActiveOrNot }">
        <span>{{ icon }}</span>
        <h3 v-if="expanded">{{ title }}</h3>
    </router-link>
</template>

<style scoped lang="scss">
.sidebartab {
    display: flex;
    align-items: center;
    margin: var(--font-s) 0;
    border-left: solid 4px transparent;
    max-height: var(--font-xl);
    padding-right: var(--font-s);
    transition: background-color 0.2s;
    border-radius: 0 5px 5px 0;
    white-space: nowrap;

    &:hover {
        background-color: var(--color-component-background);
    }

    span {
        font-family: var(--font-icon);
        margin-left: var(--font-s);
        font-size: calc(var(--font-m) * var(--font-to-icon));
    }

    h3 {
        font-size: var(--font-m);
        margin-left: var(--font-s);
        font-weight: 400;
    }


}

.link-active {
    border-left-color: var(--color-highlight);
    background-color: var(--color-component-background);

    span,
    h3 {
        color: var(--color-highlight);
    }
}
</style>