<template>
  <div>
    <h2>Login</h2>
    <form @submit.prevent="login">
      <label>
        Usuário:
        <input v-model="username" required />
      </label>
      <br />
      <label>
        Senha:
        <input type="password" v-model="password" required />
      </label>
      <br />
      <button type="submit">Entrar</button>
    </form>
    <p v-if="error" style="color: red">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')

function login() {
  if (username.value === 'admin' && password.value === '1234') {
    router.push('/')
  } else {
    error.value = 'Usuário ou senha incorretos.'
  }
}
</script>

