<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../store/authStore'
const VITE_APP_TITLE = import.meta.env.VITE_APP_TITLE

const email = ref('tuic-admin@gov.taipei')
const password = ref('vrcwzaK200Kqq9UC6j9EBbYF3YqEDc')
const authStore = useAuthStore()
</script>

<template>
    <div class="login">

        <div class="login-cont">
            <div class="login-cont-logo">
                <img src="../assets/images/TUIC.svg" alt="tuic logo" />
                <div>
                    <h1>{{ VITE_APP_TITLE }}</h1>
                    <h2>Taipei City Dashboard</h2>
                </div>
            </div>
            <form>
                <label for="email">信箱
                    <input name="email" type="email" v-model="email" /></label>
                <label for="password">密碼
                    <input name="password" type="password" v-model="password" /></label>
                <button @click.prevent="authStore.handleEmailLogin(email, password)">登入</button>
            </form>
        </div>
    </div>
</template>

<style scoped lang="scss">
.login {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    background: linear-gradient(rgba(0, 0, 0, 0.6), rgba(0, 0, 0, 0.4)), url(../assets/images/2-days-in-Taipei-Taiwan-6.jpg) center no-repeat;
    // background-image: ;
    // background-repeat: no-repeat;
    background-size: cover;
    background-position: 50%;

    &-cont {
        padding: var(--font-l);
        background-color: rgba(0, 0, 0, 0.6);
        border-radius: 5px;
        box-shadow: 0 5px 5px rgba(0, 0, 0, .5);

        &-logo {
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 2rem;

            h1 {
                font-weight: 500;
                font-size: 1rem;
            }

            h2 {
                font-size: 0.75rem;
                font-weight: 400;
            }

            img {
                filter: invert(1);
                height: 30px;
                margin: 0 var(--font-s)
            }
        }

        form {
            display: flex;
            flex-direction: column;
            width: 300px;

            label {
                display: flex;
                align-items: center;
                margin-bottom: var(--font-s);

                input {
                    border-radius: 5px;
                    border: none;
                    margin-left: var(--font-s);
                    font-size: var(--font-m);
                    color: black;
                    padding: 4px 8px;
                    flex-grow: 1;
                    background-color: white;
                }

                input:focus {
                    outline: none;
                }
            }

            button {
                background-color: var(--color-border);
                width: 5rem;
                align-self: center;
                padding: 5px 0;
                border-radius: 2rem;
                color: white;
                opacity: 1;
                transition: background-color 0.2s;
                margin-top: 1.5rem;

                &:hover {
                    background-color: var(--color-highlight);
                }
            }
        }

    }
}
</style>