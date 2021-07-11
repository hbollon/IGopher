<template>
    <nav class="navbar navbar-dark align-items-start sidebar sidebar-dark accordion bg-gradient-primary p-0">
        <div class="container-fluid d-flex flex-column p-0"><a class="navbar-brand d-flex justify-content-center align-items-center sidebar-brand m-0" href="#">
                <div class="sidebar-brand-icon rotate-n-15"><i class="fas fa-robot"></i></div>
                <div class="sidebar-brand-text mx-3"><span>Igopher</span></div>
            </a>
            <hr class="sidebar-divider my-0">
            <ul class="navbar-nav text-light" id="accordionSidebar">
                <li class="nav-item"></li>
                <li class="nav-item"></li>
                <li class="nav-item"><router-link class="nav-link" to="/"><i class="fa fa-envelope"></i><span>DM Automation&nbsp;</span></router-link></li>
                <li class="nav-item"></li>
                <li class="nav-item"><router-link class="nav-link" to="/settings"><i class="fa fa-gears"></i><span>Settings</span></router-link></li>
                <li class="nav-item"><router-link class="nav-link" to="/logs"><i class="fa fa-list-alt"></i><span>Logs</span></router-link></li>
            </ul>
            <div class="text-center d-none d-md-inline"><button class="btn rounded-circle border-0" id="sidebarToggle" type="button"></button></div>
        </div>
    </nav>
</template>

<script lang="ts">
import { Vue } from 'vue-class-component';

export default class LateralNav extends Vue {}
</script>