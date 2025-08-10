<script setup>
import Header from './components/Header.vue'
import Footer from "@/components/Footer.vue";
import {useAuthStore} from "@/stores/useAuth.js";
import {Toaster} from '@/components/ui/sonner'
import 'vue-sonner/style.css'
import {onMounted, watch} from "vue";
import {useWishlistStore} from "@/stores/useWishlist.js";

const {fetchUser, isAuthenticated} = useAuthStore()
const {initFromLocalStorage} = useWishlistStore()

onMounted(async () => {
  await fetchUser()
  initFromLocalStorage()
})

watch(isAuthenticated, async () => {
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
