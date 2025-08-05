<script setup>
import Header from './components/Header.vue'
import Footer from "@/components/Footer.vue";
import {useAuthStore} from "@/stores/useAuth.js";
import {Toaster} from '@/components/ui/sonner'
import 'vue-sonner/style.css'
import {computed, onMounted} from "vue";

const {fetchUser, user, error, isAuthenticated} = useAuthStore()

const upToDateUser = computed(() => user)
const upToDateError = computed(() => error)
const upToDateIsAuthenticated = computed(() => isAuthenticated)

console.log(upToDateUser.value)
console.log(upToDateError.value)
console.log(upToDateIsAuthenticated.value)

onMounted(async () => {
  await fetchUser()
})
</script>

<template>
  <Toaster/>
  <div class="w-full">
    <Header/>
    <div class="mt-[68px]">
      <router-view></router-view>
    </div>
    <Footer/>
  </div>
</template>
