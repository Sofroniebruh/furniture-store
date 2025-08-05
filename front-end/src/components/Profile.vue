<script setup>

import {LogOut} from "lucide-vue-next";
import {Button} from "@/components/ui/button/index.js";
import {useAuthStore} from "@/stores/useAuth.js";
import ProfileTabs from "@/components/ProfileTabs.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";

const auth = useAuthStore()
const screenSheet = useScreenSheetStore()

const handleLogout = () => {
  auth.logout()
  screenSheet.setOpenSheet(false, {name: 'ScreenProfile'})
  screenSheet.setAllDialogsClosed()
}
</script>

<template>
  <div class="flex flex-col h-full w-full px-4" v-if="auth.isAuthenticated">
    <div>
      <Button @click.stop="handleLogout" class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">Log out
        <LogOut/>
      </Button>
    </div>
    <div class="flex flex-col justify-center items-center gap-3 my-16">
      <h1 class="font-semibold text-3xl">Welcome back,</h1>
      <span class="font-base text-lg">{{ auth.user.email }}</span>
    </div>
    <ProfileTabs></ProfileTabs>
  </div>
</template>