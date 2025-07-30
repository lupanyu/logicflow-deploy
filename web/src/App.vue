<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation Header -->
    <header class="bg-white shadow">
      <div class="container mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <ServerIcon class="h-6 w-6 text-primary" />
            <h1 class="ml-2 text-xl font-bold text-gray-800">Deployment Portal</h1>
          </div>
          <nav class="flex space-x-4">
            <button 
            @click="switchTab('home')" 
              :class="['px-3 py-2 rounded-md text-sm font-medium', 
                activeTab === 'home' ? 'bg-primary text-white' : 'text-gray-600 hover:bg-gray-100']">
              <HomeIcon class="h-4 w-4 inline mr-1" />
              Home
            </button>
            <button 
            @click="switchTab('templates')" 
              :class="['px-3 py-2 rounded-md text-sm font-medium', 
                activeTab === 'templates' ? 'bg-primary text-white' : 'text-gray-600 hover:bg-gray-100']">
              <LayoutTemplateIcon class="h-4 w-4 inline mr-1" />
              Templates
            </button>
            <button 
            @click="switchTab('history')" 
              :class="['px-3 py-2 rounded-md text-sm font-medium', 
                activeTab === 'history' ? 'bg-primary text-white' : 'text-gray-600 hover:bg-gray-100']">
              <HistoryIcon class="h-4 w-4 inline mr-1" />
              History
            </button>
            <!-- <button 
            @click="switchTab('working')" 
              :class="['px-3 py-2 rounded-md text-sm font-medium', 
                activeTab === 'working' ? 'bg-primary text-white' : 'text-gray-600 hover:bg-gray-100']">
              <GitBranchIcon class="h-4 w-4 inline mr-1" />
              Working
            </button> -->
          </nav>
        </div>
      </div>
    </header>
  
    <!-- Main Content -->
    <main class="container mx-auto px-4 py-6">
      <!-- New Flow Template Creation -->
      <div v-if="showNewFlow" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <div class="flex items-center">
              <button 
                @click="showNewFlow = false" 
                class="mr-3 text-gray-500 hover:text-gray-700">
                <ArrowLeftIcon class="h-5 w-5" />
              </button>
              <h2 class="text-2xl font-bold text-gray-800">New Template</h2>
            </div>
            <div class="flex space-x-3">
              <button 
                @click="handleNewFlowSave()"
                class="bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center">
                <CheckCircleIcon class="h-4 w-4 mr-1" />
                Save Template
              </button>
            </div>
          </div>
        </div>
        <div class="p-6">
          <div class="mb-6">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Template Name</label>
              <input 
                v-model="newTemplateData.name" 
                type="text" 
                placeholder="Enter template name" 
                class="w-full border-gray-300 rounded-md shadow-sm focus:border-primary focus:ring focus:ring-primary focus:ring-opacity-50"
              />
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Template Env</label>
              <select 
                v-model="newTemplateData.env" 
                class="w-full border-gray-300 rounded-md shadow-sm focus:border-primary focus:ring focus:ring-primary focus:ring-opacity-50">
                <option value="Test">Test</option>
                <option value="Pre">Pre</option>
                <option value="Prod">Prod</option>
              </select>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea 
                v-model="newTemplateData.description" 
                placeholder="Enter template description" 
                class="w-full border-gray-300 rounded-md shadow-sm focus:border-primary focus:ring focus:ring-primary focus:ring-opacity-50"
                rows="3"
              ></textarea>
            </div>
          </div>
          
          <!-- LogicFlow Canvas for Template Creation with Save Button -->
          <div class="border rounded-lg relative" style="height: 500px;">
            <NewFlow ref="newFlowRef" @save="handleNewFlowSave" />--------------------
          </div>
        </div>
      </div>
  
      <!-- Template Details with LogicFlow Canvas -->
      <div v-else-if="selectedTemplate" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <div class="flex items-center">
              <button 
                @click="selectedTemplate = null" 
                class="mr-3 text-gray-500 hover:text-gray-700">
                <ArrowLeftIcon class="h-5 w-5" />
              </button>
              <h2 class="text-2xl font-bold text-gray-800">{{ selectedTemplate.name }}</h2>
            </div>
            <div class="flex space-x-3">
              <button 
                @click="executeDeployment(selectedTemplate)" 
                class="bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center"
                :disabled="isDeploying">
                <RocketIcon class="h-4 w-4 mr-1" />
                {{ isDeploying ? 'Deploying...' : 'Deploy' }}
              </button>
              <button 
                @click="handleUpdateTemplate()" 
                class="bg-orange-600 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center">
                <CheckCircleIcon class="h-4 w-4 mr-1" />
                Save Template
              </button>
            </div>
          </div>
        </div>
        <div class="p-6">
          <div class="mb-4">
            <span :class="['text-xs px-2 py-1 rounded-full', 
              selectedTemplate.type === 'Frontend' ? 'bg-blue-100 text-blue-800' : 
              selectedTemplate.type === 'Backend' ? 'bg-green-100 text-green-800' : 
              'bg-purple-100 text-purple-800']">
              {{ selectedTemplate.type }}
            </span>
            <p class="mt-4 text-gray-700">{{ selectedTemplate.description }}</p>
          </div>
          
          <!-- LogicFlow Canvas for Template -->
          <div class="border rounded-lg mt-4" style="height: 500px;">
            <LF :flow-name="selectedTemplate.name" :readonly="true" ref="LFRef" />
          </div>
        </div>
      </div>
  
      <!-- Deployment Details with ExecutionDetail -->
      <div v-else-if="selectedDeployment" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <div class="flex items-center">
              <button 
                @click="selectedDeployment = null" 
                class="mr-3 text-gray-500 hover:text-gray-700">
                <ArrowLeftIcon class="h-5 w-5" />
              </button>
              <h2 class="text-2xl font-bold text-gray-800">Deployment: {{ selectedDeployment.name || selectedDeployment.flowId }}</h2>
            </div>
            <div class="flex space-x-3">
              <!-- <button 
                v-if="selectedDeployment.status !== 'In Progress'"
                @click="redeployFromHistory(selectedDeployment)"
                class="bg-primary text-white px-4 py-2 rounded-md text-sm font-medium flex items-center"
                :disabled="isDeploying">
                <RocketIcon class="h-4 w-4 mr-1" />
                {{ isDeploying ? 'Redeploying...' : 'Redeploy' }}
              </button> -->
            </div>
          </div>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-2 gap-4 mb-6" >
            <div>
              <p class="text-sm text-gray-500">Status</p>
              <span :class="['px-2 py-1 text-xs rounded-full', 
                selectedDeployment?.status === 'success' ? 'bg-green-100 text-green-800' : 
                selectedDeployment?.status === 'failed' ? 'bg-red-100 text-red-800' : 
                'bg-yellow-100 text-yellow-800']">
                {{  selectedDeployment.status }}
              </span>
            </div>
            <div>
              <p class="text-sm text-gray-500">Started At</p>
              <p class="font-medium">{{ selectedDeployment.startTime }}</p>
            </div>

            <div v-if="  selectedDeployment.env">
              <p class="text-sm text-gray-500">Env</p>
              <p class="font-medium">{{  selectedDeployment.env }}</p>
            </div>
            <div v-if=" selectedDeployment.endTime">
              <p class="text-sm text-gray-500">Ended At</p>
              <p class="font-medium">{{ selectedDeployment.endTime }}</p>
            </div>
            <div v-if="selectedDeployment.duration">
              <p class="text-sm text-gray-500">SpendTime</p>
              <p class="font-medium">{{selectedDeployment.duration }}</p>
            </div>
            <div>
              <p>
             <StopDeploymentButton
                      v-if="selectedDeployment?.status === 'running'"
                      :deployment-id="selectedDeployment?.flowId"
                      @success="handleStopSuccess"
                      size="small"
                    />
              </p>
            </div>
          </div>
          
          <!-- ExecutionDetail Component for Deployment -->
          <ExecutionDetail :getSelectedDeployment="getSelectedDeployment"
            :setSelectedDeployment="setSelectedDeployment"
            ref="ExecutionDetailRef"
           />
        </div>
      </div>
  
      <!-- Home View -->
      <div v-else-if="activeTab === 'home'" class="bg-white rounded-lg shadow p-6">
        <h2 class="text-2xl font-bold text-gray-800 mb-4">Welcome to Deployment Portal</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-blue-50 p-4 rounded-lg border border-blue-100">
            <LayoutTemplateIcon class="h-8 w-8 text-blue-500 mb-2" />
            <h3 class="text-lg font-semibold mb-2">Templates</h3>
            <p class="text-gray-600">Browse and select from saved deployment templates.</p>
            <button 
            @click="switchTab('templates')" 
              class="mt-4 text-blue-600 hover:text-blue-800 font-medium flex items-center">
              View Templates
              <ArrowRightIcon class="h-4 w-4 ml-1" />
            </button>
          </div>
          <div class="bg-green-50 p-4 rounded-lg border border-green-100">
            <HistoryIcon class="h-8 w-8 text-green-500 mb-2" />
            <h3 class="text-lg font-semibold mb-2">History</h3>
            <p class="text-gray-600">View your past deployments and their status.</p>
            <button 
            @click="switchTab('history')" 
              class="mt-4 text-green-600 hover:text-green-800 font-medium flex items-center">
              View History
              <ArrowRightIcon class="h-4 w-4 ml-1" />
            </button>
          </div>
          <div class="bg-purple-50 p-4 rounded-lg border border-purple-100">
            <RocketIcon class="h-8 w-8 text-purple-500 mb-2" />
            <h3 class="text-lg font-semibold mb-2">Working</h3>
            <p class="text-gray-600">View the current deployment in working</p>
            <button 
            @click="switchTab('working')" 
              class="mt-4 text-purple-600 hover:text-purple-800 font-medium flex items-center">
              View Working
              <ArrowRightIcon class="h-4 w-4 ml-1" />
            </button>
          </div>
        </div>
      </div>
  
      <!-- Working View -->
      <div v-else-if="activeTab === 'working'" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <h2 class="text-2xl font-bold text-gray-800">Working Pipelines</h2>
          </div>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div 
              v-for="pipeline in pipelines" 
              :key="pipeline.id" 
              class="border rounded-lg p-4 hover:shadow-md cursor-pointer transition-shadow"
              @click="viewWorkingDetails(pipeline)">
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="font-semibold text-lg">{{ pipeline.name }}</h3>
                  <p class="text-gray-600 text-sm mt-1">{{ pipeline.description }}</p>
                </div>
                <span :class="['text-xs px-2 py-1 rounded-full', 
                  pipeline.category === 'CI/CD' ? 'bg-blue-100 text-blue-800' : 
                  pipeline.category === 'Data' ? 'bg-green-100 text-green-800' : 
                  'bg-purple-100 text-purple-800']">
                  {{ pipeline.category }}
                </span>
              </div>
              <div class="mt-4 flex justify-between items-center">
                <span class="text-sm text-gray-500">Last updated: {{ pipeline.updatedAt }}</span>
                <div class="flex space-x-2">
                  <button class="text-primary hover:text-primary-dark">
                    <ExternalLinkIcon class="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
  
      <!-- Templates View -->
      <div v-else-if="activeTab === 'templates'" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <h2 class="text-2xl font-bold text-gray-800">Templates List</h2>
            <div class="flex space-x-3">
              <button 
                @click="openNewTemplateFlow()" 
                class="bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center">
                <PlusIcon class="h-4 w-4 mr-1" />
                New Template
              </button>
            </div>
          </div>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div 
              v-for="template in templates" 
              :key="template.id" 
              class="border rounded-lg p-4 hover:shadow-md cursor-pointer transition-shadow"
              @click="viewTemplateDetails(template)">
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="font-semibold text-lg">{{ template.name }}</h3>
                </div>
              </div>
              <div class="mt-4 flex justify-between items-center">
                <span class="text-sm text-gray-500">Last updated: {{ template.updatedAt }}</span>
                <div class="flex space-x-2">
                  <button 
                    @click.stop="executeDeployment(template.name)" 
                    class="text-green-600 hover:text-green-800 flex items-center"
                    :disabled="isDeploying">
                    <RocketIcon class="h-4 w-4 mr-1" />
                    Deploy
                  </button>
                  <button class="text-primary hover:text-primary-dark">
                    <ExternalLinkIcon class="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
  
      <!-- History View -->
      <div v-else-if="activeTab === 'history'" class="bg-white rounded-lg shadow">
        <div class="p-6 border-b border-gray-200">
          <h2 class="text-2xl font-bold text-gray-800">History</h2>
        </div>
        <div class="p-6">
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Start Time</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">End Time</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">SpendTime</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="deployment in sortedDeployHistory" :key="deployment.flowId" class="hover:bg-gray-50">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ deployment.flowId }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ deployment.name }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span :class="['px-2 py-1 text-xs rounded-full', 
                      deployment.status === 'success' ? 'bg-green-100 text-green-800' : 
                      deployment.status === 'failed' ? 'bg-red-100 text-red-800' : 
                      deployment.status === 'running' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-yellow-100 text-yellow-800']">
                      {{ deployment.status }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ deployment.startTime }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ deployment.endTime }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ deployment.duration }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button 
                      @click="viewDeploymentDetails(deployment)" 
                      class="text-primary hover:text-primary-dark mr-3">
                      Details
                    </button>
                    <StopDeploymentButton
                      v-if="deployment.status === 'running'"
                      :deployment-id="deployment.flowId"
                      @success="handleStopSuccess"
                      size="small"
                    />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
  
      <!-- Deployment Status Toast -->
      <div 
        v-if="deploymentToast.show" 
        :class="['fixed bottom-4 right-4 p-4 rounded-md shadow-lg max-w-md z-50 flex items-center', 
          deploymentToast.type === 'success' ? 'bg-green-100 border-l-4 border-green-500' : 
          deploymentToast.type === 'error' ? 'bg-red-100 border-l-4 border-red-500' : 
          'bg-blue-100 border-l-4 border-blue-500']">
        <div :class="['mr-3', 
          deploymentToast.type === 'success' ? 'text-green-500' : 
          deploymentToast.type === 'error' ? 'text-red-500' : 
          'text-blue-500']">
          <CheckCircleIcon v-if="deploymentToast.type === 'success'" class="h-5 w-5" />
          <AlertCircleIcon v-if="deploymentToast.type === 'error'" class="h-5 w-5" />
          <InfoIcon v-if="deploymentToast.type === 'info'" class="h-5 w-5" />
        </div>
        <div>
          <p :class="['font-medium', 
            deploymentToast.type === 'success' ? 'text-green-800' : 
            deploymentToast.type === 'error' ? 'text-red-800' : 
            'text-blue-800']">
            {{ deploymentToast.title }}
          </p>
          <p :class="['text-sm', 
            deploymentToast.type === 'success' ? 'text-green-700' : 
            deploymentToast.type === 'error' ? 'text-red-700' : 
            'text-blue-700']">
            {{ deploymentToast.message }}
          </p>
        </div>
        <button 
          @click="deploymentToast.show = false" 
          class="ml-auto text-gray-500 hover:text-gray-700">
          <XIcon class="h-4 w-4" />
        </button>
      </div>
    </main>
  </div>
  </template>
  
  <script setup>
  import { ref, onMounted, computed } from 'vue';
  import { ElMessageBox } from 'element-plus';
  import { 
    ServerIcon, 
    HomeIcon, 
    LayoutTemplateIcon, 
    HistoryIcon, 
    ArrowRightIcon, 
    PlusIcon, 
    ExternalLinkIcon,
    XIcon,
    RocketIcon,
    CheckCircleIcon,
    AlertCircleIcon,
    InfoIcon,
    GitBranchIcon,
    ArrowLeftIcon,
  } from 'lucide-vue-next';
  
  // Import your existing LogicFlow components
  import ExecutionDetail from './components/ExecutionDetail.vue';
  import LF from './components/LF.vue';
  import NewFlow from './components/NewFlow.vue';
  import StopDeploymentButton from  './components/StopDeploymentButton.vue';
  // View states
  const activeTab = ref('home');
  const selectedTemplate = ref(null);
  const selectedDeployment = ref(null);
  const showNewFlow = ref(false);
  const newFlowRef = ref(null);
  const LFRef = ref(null);
  const ExecutionDetailRef = ref(null);
  // New template data
  const newTemplateData = ref({
    name: '',
    env: 'Test',
    description: '',
    nodes: null,
    edges: null
  });
  
  const updateTemplateData = ref({
    name: '',
    env: '',
    description: '',
    nodes: null,
    edges: null
  });
  
  // Deployment state
  const isDeploying = ref(false);
   
  const deploymentToast = ref({
    show: false,
    type: 'info',
    title: '',
    message: '',
    timeout: null
  });
  
  // Function to open new template flow
  const openNewTemplateFlow = () => {
    activeTab.value = 'templates';
    newTemplateData.value = {
      name: '',
      env: 'Test',
      description: '',
      nodes: null,
      edges: null
    };
    showNewFlow.value = true;
  };
  
  const getSelectedDeployment = () => {
    return selectedDeployment.value;
  };

  const setSelectedDeployment = (deployment) => {
    selectedDeployment.value = deployment;
  }

  
  // Function to handle new flow save from NewFlow component
  const handleNewFlowSave = () => {
    const flowData = newFlowRef.value.GetGraphData();
    newTemplateData.value.name = newTemplateData.value.name.trim();
    newTemplateData.value.description = newTemplateData.value.description.trim();
    newTemplateData.value.env = newTemplateData.value.env.trim();
    newTemplateData.value.nodes = flowData.nodes;
    newTemplateData.value.edges = flowData.edges;
    console.log('newTemplateData.value', newTemplateData.value);
    saveNewTemplate();
  };
  
  const handleUpdateTemplate = () => {
    const flowData = LFRef.value.LFGetGraphData();
    console.log(flowData);
  
    updateTemplateData.value.name = selectedTemplate.value.name;
    updateTemplateData.value.env = selectedTemplate.value.env;
    updateTemplateData.value.description = selectedTemplate.value.description;
    updateTemplateData.value.nodes = flowData.nodes;
    updateTemplateData.value.edges = flowData.edges;
    console.log('updateTemplateData.value', updateTemplateData.value);
    updateTemplate();
  };
  
  const updateTemplate = async () => {
    try {
      const response = await fetch("/api/v1/flow/"+updateTemplateData.value.name, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateTemplateData.value)
      });
      console.log(response);
      if (!response.ok) {
        showToast(
          'error',
          'Update Failed',
          'Update template failed:'+response.statusText+response.body
        );
        return;
      }
      showToast(
        'success',
        'Update Success',
        'Update template success'
      );
    } catch (e) {
      console.log('保存template失败:', e);
      showToast(
        'error',
        'Update Failed',
        'Update template failed: ' + e.message
      );
    }
  };
  
  // Function to save new template
  const saveNewTemplate = async () => {
    if (!newTemplateData.value.name) {
      showToast('error', 'Validation Error', 'Template name is required');
      return;
    }

    console.log('newTemplateData.value:', newTemplateData.value);
    try {
      const response = await fetch("/api/v1/flow/", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newTemplateData.value)
      });
      console.log(response);
      if (!response.ok) {
        showToast(
          'error',
          'Save Failed',
          'Save template failed:'+response.statusText+response.body
        );
        return;
      }
      showToast(
        'success',
        'Save Success',
        'Save template success'
      );
      // Refresh templates list
      await fetchFlows();
    } catch (e) {
      console.log('保存template失败:', e);
      showToast(
        'error',
        'Save Failed',
        'Save template failed: ' + e.message
      );
    }
  };
  
  const executeDeployment = async (template) => {
                // 新增确认弹窗

    try {
      await ElMessageBox.confirm(
                'Are you sure to deploy this template?', 
                'Confirm Deployment', 
                {
                    confirmButtonText: 'Confirm',
                    cancelButtonText: 'Cancel',
                    type: 'warning'
                }
            )
            
    isDeploying.value = true;
    console.log("template", template);
      let url = '/api/v1/deploy/';
      let args = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({})
      };
      if (typeof template === 'string') {
        url = '/api/v1/deploy/' + template;
       }else{
        url = '/api/v1/deploy/' + template.name;
       }
      console.log("url", url);
      console.log("args", args);
  
      const response = await fetch(url, args);
      console.log("response", response);
      const data = await response.json();
      if (!response.ok) {
        showToast('error', 'Deployment Failed', data.error || 'There was an error starting the deployment. Please try again.');
        return;
      }
      // Show success toast
      showToast('success', 'Deployment Started', (selectedTemplate.value?.name || template) + ' is now being deployed ...');
      
      // Switch to history tab to show the deployment progress
      switchTab('history');
      
    } catch (error) {
      // Show error toast
      showToast('error', 'Deployment Failed', error.message || 'There was an error starting the deployment. Please try again.');
      console.error('Deployment error:', error);
    } finally {
      isDeploying.value = false;
    }
  };
  
  const showToast = (type, title, message) => {
    // Clear any existing timeout
    if (deploymentToast.value.timeout) {
      clearTimeout(deploymentToast.value.timeout);
    }
  
    // Set toast data
    deploymentToast.value = {
      show: true,
      type,
      title,
      message,
      timeout: setTimeout(() => {
        deploymentToast.value.show = false;
      }, 5000) // Hide after 5 seconds
    };
  };
  
  // Functions to view details
  const viewTemplateDetails = (template) => {
    selectedTemplate.value = template;
  };
  
  const viewDeploymentDetails = (deployment) => {
    selectedDeployment.value = deployment;
  };
  
  // Function to view working details
  const viewWorkingDetails = (pipeline) => {
    // Create a deployment-like object from the pipeline for ExecutionDetail
    const workingDeployment = {
      flowId: pipeline.id,
      name: pipeline.name,
      status: 'In Progress',
      startTime: pipeline.updatedAt,
      description: pipeline.description,
      // Add any other properties needed by ExecutionDetail
    };
  
    selectedDeployment.value = workingDeployment;
  };
  
  const templates = ref([]);
  
  // Function to fetch flows
  async function fetchFlows() {
    try {
      const response = await fetch('/api/v1/flow', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      if (!response.ok) throw new Error('Failed to fetch data');
      const result = await response.json();
      console.log(result);
      templates.value = result.data || [];
    } catch (e) {
      showToast('error', 'Error get flowlist', e.message);
    } 
  }
  
  // Mock data for pipelines
  const pipelines = ref([
    {
      id: 1,
      name: 'Frontend Deployment',
      description: 'CI/CD pipeline for React applications',
      category: 'CI/CD',
      updatedAt: '2023-10-20'
    },
    {
      id: 2,
      name: 'Data Processing',
      description: 'ETL pipeline for data transformation',
      category: 'Data',
      updatedAt: '2023-10-18'
    },
    {
      id: 3,
      name: 'Microservice Deployment',
      description: 'Deployment pipeline for microservices',
      category: 'CI/CD',
      updatedAt: '2023-10-15'
    },
    {
      id: 4,
      name: 'Database Migration',
      description: 'Pipeline for database schema migrations',
      category: 'Data',
      updatedAt: '2023-10-12'
    },
    {
      id: 5,
      name: 'Serverless Deployment',
      description: 'Pipeline for serverless function deployment',
      category: 'Serverless',
      updatedAt: '2023-10-05'
    }
  ]);
  
  const deployHistory = ref([]);
  
  // Computed property to sort deployment history by start time (newest first)
  const sortedDeployHistory = computed(() => {
    return [...deployHistory.value].sort((a, b) => {
      // Convert dates to timestamps for comparison
      const dateA = new Date(a.startTime || 0).getTime();
      const dateB = new Date(b.startTime || 0).getTime();
      // Sort in descending order (newest first)
      return dateB - dateA;
    });
  });
  
  // Function to fetch deployments
  const fetchDeployHistory = async () => {
    try {
      const response = await fetch('/api/v1/deploy', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      if (!response.ok) throw new Error('Failed to fetch data');
      const result = await response.json();
      console.log(result);
      deployHistory.value = result || [];
    } catch (e) {
      showToast('error', 'Error get deploy history', e.message);
    } 
  };
  
  // Function to switch tabs and fetch relevant data
  const switchTab = async (tab) => {
    activeTab.value = tab;
    selectedTemplate.value = null;
    selectedDeployment.value = null;
    showNewFlow.value = false;
  
    // Fetch data based on the selected tab
    if (tab === 'templates') {
      await fetchFlows();
    } else if (tab === 'history') {
      await fetchDeployHistory();
    } else if (tab === 'working') {
      await fetchWorkingPipelines();
    }
  };
  
  // Function to fetch working pipelines
  const fetchWorkingPipelines = async () => {
    try {
      console.log('Fetching working pipelines');
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300));
    } catch (e) {
      showToast('error', 'Error fetching working pipelines', e.message);
    }
  };
  
  onMounted(() => {
    // Initial app setup if needed
  });
  </script>
  
  <style>
  :root {
    --color-primary: #4f46e5;
    --color-primary-dark: #4338ca;
  }
  
  .bg-primary {
    background-color: var(--color-primary);
  }
  
  .bg-primary-dark {
    background-color: var(--color-primary-dark);
  }
  
  .text-primary {
    color: var(--color-primary);
  }
  
  .text-primary-dark {
    color: var(--color-primary-dark);
  }
  
  .hover\:bg-primary-dark:hover {
    background-color: var(--color-primary-dark);
  }
  
  .hover\:text-primary-dark:hover {
    color: var(--color-primary-dark);
  }
  </style>