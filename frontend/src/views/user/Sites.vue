<template>
  <div class="sites-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>站点查询</span>
          <div class="search-box">
            <el-input
              v-model="keyword"
              placeholder="搜索站点名称或地址"
              style="width: 200px"
              clearable
              @keyup.enter="loadSites"
            >
              <template #append>
                <el-button @click="loadSites">
                  <el-icon><Search /></el-icon>
                </el-button>
              </template>
            </el-input>
          </div>
        </div>
      </template>
      <el-table :data="sites" style="width: 100%">
        <el-table-column prop="name" label="站点名称" min-width="180" />
        <el-table-column prop="province" label="省份" width="100" />
        <el-table-column prop="city" label="城市" width="100" />
        <el-table-column prop="district" label="区县" width="100" />
        <el-table-column prop="address" label="详细地址" min-width="250" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="manager" label="负责人" width="100" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '正常' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSitesAPI } from '@/api'

const keyword = ref('')
const sites = ref<any[]>([])

const loadSites = async () => {
  try {
    const params = keyword.value ? { keyword: keyword.value } : undefined
    const res: any = await getSitesAPI(params)
    sites.value = res.data || []
  } catch (error) {
    console.error('加载站点失败:', error)
  }
}

onMounted(() => {
  loadSites()
})
</script>

<style scoped>
.sites-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
