<script setup>
const props = defineProps({ for: String })
</script>

<template>
    <label :for="props.for" class="check">
        <svg width="16px" height="16px" viewBox="0 0 18 18">
            <path
                d="M1,9 L1,3.5 C1,2 2,1 3.5,1 L14.5,1 C16,1 17,2 17,3.5 L17,14.5 C17,16 16,17 14.5,17 L3.5,17 C2,17 1,16 1,14.5 L1,9 Z">
            </path>
            <polyline points="1 9 7 14 15 4"></polyline>
        </svg>
        <slot></slot>
    </label>
</template>

<style scoped lang="scss">
.customcheckbox {
    display: flex;
    flex-direction: row;
}

.check {
    cursor: pointer;
    position: relative;
    width: 16px;
    height: 16px;
    -webkit-tap-highlight-color: transparent;
    transform: translate3d(0, 0, 0);
}

svg {
    margin-right: 0.5rem;
}

// .check:before {
//     content: "";
//     position: absolute;
//     top: -15px;
//     left: -15px;
//     width: 48px;
//     height: 48px;
//     border-radius: 50%;
//     background: rgba(34, 50, 84, 0.03);
//     opacity: 0;
//     transition: opacity 0.2s ease;
// }

.check svg {
    position: relative;
    z-index: 1;
    fill: none;
    stroke-linecap: round;
    stroke-linejoin: round;
    stroke: #c8ccd4;
    stroke-width: 1.5;
    transform: translate3d(0, 2px, 0);
    transition: all 0.2s ease;
}

.check svg path {
    stroke-dasharray: 60;
    stroke-dashoffset: 0;
}

.check svg polyline {
    stroke-dasharray: 22;
    stroke-dashoffset: 66;
}

.check:hover:before {
    opacity: 1;
}

.check:hover svg {
    stroke: #4285f4;
}

.custom-check-input:checked+.check svg {
    stroke: #4285f4;
}

.custom-check-input:checked+.check svg path {
    stroke-dashoffset: 60;
    transition: all 0.3s linear;
}

.custom-check-input:checked+.check svg polyline {
    stroke-dashoffset: 42;
    transition: all 0.2s linear;
    transition-delay: 0.15s;
}
</style>