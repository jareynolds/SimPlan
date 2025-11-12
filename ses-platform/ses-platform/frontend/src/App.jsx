import React, { useState, useEffect } from 'react';
import { AlertCircle, Check, ChevronRight, ChevronLeft, Upload, Play, Pause, Trash2, Search, Filter, DollarSign, Clock, Settings, Database, Network, Server, Box } from 'lucide-react';

// Capability and Enabler data from markdown files
const CAPABILITIES = [
  { id: 'C01', name: 'Spec Authoring & Validation', description: 'Web editor and services to author SES with validation', enablers: ['E02', 'E17', 'E11', 'E20'], dependencies: [] },
  { id: 'C02', name: 'Parsing & Internal Modeling', description: 'Convert validated SES into normalized internal models', enablers: ['E03', 'E04', 'E17', 'E01'], dependencies: ['C01'] },
  { id: 'C03', name: 'Planning Engine', description: 'Generates execution plan with resource allocation', enablers: ['E03', 'E04', 'E08', 'E16'], dependencies: ['C01', 'C02'] },
  { id: 'C04', name: 'Provisioning Automation', description: 'Executes plan to create infrastructure', enablers: ['E05', 'E04', 'E01', 'E18', 'E20'], dependencies: ['C03', 'C09'] },
  { id: 'C05', name: 'Orchestration & State Machine', description: 'Manages lifecycle phases and state transitions', enablers: ['E04', 'E03', 'E18', 'E01'], dependencies: ['C04', 'C17'] },
  { id: 'C06', name: 'Monitoring & Metrics', description: 'Collects infrastructure and system metrics', enablers: ['E06', 'E07', 'E11', 'E19'], dependencies: ['C04', 'C12'] },
  { id: 'C07', name: 'Logging & Audit', description: 'Structured logs with audit trail', enablers: ['E07', 'E19', 'E01', 'E09', 'E10'], dependencies: ['C16', 'C17'] },
  { id: 'C08', name: 'Cost Management', description: 'Real-time cost tracking and optimization', enablers: ['E08', 'E16', 'E06', 'E11'], dependencies: ['C06', 'C03'] },
  { id: 'C09', name: 'Security & Compliance', description: 'Credential vault, RBAC, policy engine', enablers: ['E09', 'E10', 'E19', 'E01'], dependencies: ['C07', 'C17'] },
  { id: 'C10', name: 'Spec-Kit & Reuse', description: 'Library of templates and reusable components', enablers: ['E17', 'E02', 'E11', 'E20'], dependencies: ['C17', 'C18'] },
  { id: 'C11', name: 'Reservation & Scheduling', description: 'Time-window validation and resource locking', enablers: ['E14', 'E04', 'E03', 'E09'], dependencies: ['C02', 'C17'] },
  { id: 'C12', name: 'Simulation Execution', description: 'Runs scenarios and validates assertions', enablers: ['E12', 'E13', 'E04', 'E06', 'E20'], dependencies: ['C04', 'C05'] },
  { id: 'C13', name: 'Error Handling & Recovery', description: 'Automated rollback and retry workflows', enablers: ['E15', 'E07', 'E04', 'E18', 'E13'], dependencies: ['C05', 'C07'] },
  { id: 'C14', name: 'Provider Abstraction', description: 'Unified adapter layer for cloud/on-prem', enablers: ['E05', 'E18', 'E01', 'E20'], dependencies: ['C09', 'C04'] },
  { id: 'C15', name: 'Environment Visualization', description: 'Topology and dependency dashboards', enablers: ['E11', 'E06', 'E07', 'E03', 'E08'], dependencies: ['C06', 'C08'] },
  { id: 'C16', name: 'Messaging & Agent Coordination', description: 'Internal message bus with correlation', enablers: ['E18', 'E01', 'E07', 'E09'], dependencies: ['C17', 'C09'] },
  { id: 'C17', name: 'State Persistence Layer', description: 'ACID store for environment state', enablers: ['E01', 'E19', 'E07', 'E09'], dependencies: ['C09'] },
  { id: 'C18', name: 'Access & Governance', description: 'Role management and policy enforcement', enablers: ['E09', 'E10', 'E19', 'E07', 'E11'], dependencies: ['C17', 'C07'] }
];

const ENABLERS = {
  'E01': { name: 'Core Platform Infra', description: 'DB, storage, vault, message broker' },
  'E02': { name: 'Schema & Validation', description: 'JSON/YAML validators, rule engine' },
  'E03': { name: 'Graph & Planning', description: 'DAG construction, topological sort' },
  'E04': { name: 'Execution Framework', description: 'Workflow engine, task dispatcher' },
  'E05': { name: 'Provider SDK', description: 'Cloud/on-prem SDK wrappers' },
  'E06': { name: 'Metrics Stack', description: 'Prometheus, time-series storage' },
  'E07': { name: 'Logging & Tracing', description: 'Central collector, distributed tracing' },
  'E08': { name: 'Cost Engine', description: 'Pricing cache, estimation algorithms' },
  'E09': { name: 'Security/RBAC', description: 'Auth, MFA, policy enforcement' },
  'E10': { name: 'Compliance Engine', description: 'Rule DSL, scheduled checks' },
  'E11': { name: 'UI Components', description: 'Web editor, dashboards, visualizers' },
  'E12': { name: 'Workload Drivers', description: 'Load/stress test integrations' },
  'E13': { name: 'Health & Validation', description: 'Endpoint pollers, status aggregators' },
  'E14': { name: 'Reservation Scheduler', description: 'Time-slot index, priority queue' },
  'E15': { name: 'Rollback Manager', description: 'LIFO unwinder, deallocator' },
  'E16': { name: 'Optimization Advisory', description: 'Right-sizing recommendations' },
  'E17': { name: 'Template & Catalog', description: 'Versioned spec templates' },
  'E18': { name: 'Message Protocol', description: 'Serialization, correlation IDs' },
  'E19': { name: 'Data Protection', description: 'Encryption, data classification' },
  'E20': { name: 'Testing & QA', description: 'Spec fixtures, provider mocks' }
};

// Main App Component
export default function SESPlatform() {
  const [view, setView] = useState('dashboard');
  const [environments, setEnvironments] = useState([
    {
      id: 'env-001',
      name: 'Production Test Env',
      status: 'running',
      capabilities: ['C01', 'C02', 'C03', 'C04', 'C06', 'C09'],
      cost: 145.50,
      uptime: '3d 5h',
      health: 98
    },
    {
      id: 'env-002',
      name: 'Dev Integration',
      status: 'provisioning',
      capabilities: ['C01', 'C02', 'C04', 'C12'],
      cost: 42.20,
      uptime: '12h',
      health: 85
    }
  ]);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Database className="w-8 h-8 text-blue-600" />
              <div>
                <h1 className="text-2xl font-bold text-gray-900">SES Platform</h1>
                <p className="text-sm text-gray-500">Simulation Environment Specification</p>
              </div>
            </div>
            <nav className="flex space-x-1">
              <button
                onClick={() => setView('dashboard')}
                className={`px-4 py-2 rounded-lg font-medium transition ${
                  view === 'dashboard'
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                Dashboard
              </button>
              <button
                onClick={() => setView('create')}
                className={`px-4 py-2 rounded-lg font-medium transition ${
                  view === 'create'
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                Create Environment
              </button>
              <button
                onClick={() => setView('templates')}
                className={`px-4 py-2 rounded-lg font-medium transition ${
                  view === 'templates'
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                Templates
              </button>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {view === 'dashboard' && <Dashboard environments={environments} setView={setView} />}
        {view === 'create' && <CreateEnvironmentWizard setView={setView} setEnvironments={setEnvironments} />}
        {view === 'templates' && <TemplatesLibrary />}
      </main>
    </div>
  );
}

// Dashboard Component
function Dashboard({ environments, setView }) {
  const [searchTerm, setSearchTerm] = useState('');

  const getStatusColor = (status) => {
    const colors = {
      running: 'bg-green-100 text-green-800',
      provisioning: 'bg-blue-100 text-blue-800',
      stopped: 'bg-gray-100 text-gray-800',
      error: 'bg-red-100 text-red-800'
    };
    return colors[status] || colors.stopped;
  };

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Environments</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">12</p>
            </div>
            <Server className="w-10 h-10 text-blue-600" />
          </div>
        </div>
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Cost (Today)</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">$432</p>
            </div>
            <DollarSign className="w-10 h-10 text-green-600" />
          </div>
        </div>
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Avg Health</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">94%</p>
            </div>
            <Network className="w-10 h-10 text-purple-600" />
          </div>
        </div>
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Reserved</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">3</p>
            </div>
            <Clock className="w-10 h-10 text-orange-600" />
          </div>
        </div>
      </div>

      {/* Environments List */}
      <div className="bg-white rounded-lg border border-gray-200">
        <div className="p-6 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold text-gray-900">Environments</h2>
            <button
              onClick={() => setView('create')}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-medium"
            >
              + New Environment
            </button>
          </div>
          <div className="mt-4 relative">
            <Search className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search environments..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        <div className="divide-y divide-gray-200">
          {environments.map((env) => (
            <div key={env.id} className="p-6 hover:bg-gray-50 transition">
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center space-x-3">
                    <h3 className="text-lg font-semibold text-gray-900">{env.name}</h3>
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(env.status)}`}>
                      {env.status}
                    </span>
                  </div>
                  <div className="mt-2 flex items-center space-x-6 text-sm text-gray-600">
                    <span className="flex items-center">
                      <Box className="w-4 h-4 mr-1" />
                      {env.capabilities.length} capabilities
                    </span>
                    <span className="flex items-center">
                      <DollarSign className="w-4 h-4 mr-1" />
                      ${env.cost.toFixed(2)}/day
                    </span>
                    <span className="flex items-center">
                      <Clock className="w-4 h-4 mr-1" />
                      {env.uptime}
                    </span>
                    <span>Health: {env.health}%</span>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <button className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg transition">
                    <Settings className="w-5 h-5" />
                  </button>
                  <button className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition">
                    <Play className="w-5 h-5" />
                  </button>
                  <button className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition">
                    <Upload className="w-5 h-5" />
                  </button>
                  <button className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition">
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// Create Environment Wizard
function CreateEnvironmentWizard({ setView, setEnvironments }) {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    owner: '',
    tags: '',
    capabilities: [],
    enablers: {},
    compute: { cpu: 4, memory: 16, instances: 2 },
    storage: 500,
    network: 'private',
    priority: 'medium',
    startTime: '',
    duration: 24
  });

  const totalSteps = 5;

  const updateFormData = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const toggleCapability = (capId) => {
    setFormData(prev => {
      const caps = prev.capabilities.includes(capId)
        ? prev.capabilities.filter(id => id !== capId)
        : [...prev.capabilities, capId];
      
      // Auto-select required enablers
      const requiredEnablers = {};
      caps.forEach(id => {
        const cap = CAPABILITIES.find(c => c.id === id);
        cap.enablers.forEach(enablerId => {
          requiredEnablers[enablerId] = true;
        });
      });
      
      return { ...prev, capabilities: caps, enablers: requiredEnablers };
    });
  };

  const calculateEstimatedCost = () => {
    const computeCost = formData.compute.cpu * formData.compute.memory * formData.compute.instances * 0.05;
    const storageCost = formData.storage * 0.1;
    const capabilityCost = formData.capabilities.length * 5;
    return (computeCost + storageCost + capabilityCost).toFixed(2);
  };

  const handleSubmit = () => {
    const newEnv = {
      id: `env-${Date.now()}`,
      name: formData.name,
      status: 'provisioning',
      capabilities: formData.capabilities,
      cost: parseFloat(calculateEstimatedCost()),
      uptime: '0h',
      health: 100
    };
    setEnvironments(prev => [...prev, newEnv]);
    setView('dashboard');
  };

  return (
    <div className="max-w-5xl mx-auto">
      {/* Progress Steps */}
      <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
        <div className="flex items-center justify-between">
          {[1, 2, 3, 4, 5].map((s) => (
            <React.Fragment key={s}>
              <div className="flex items-center">
                <div className={`w-10 h-10 rounded-full flex items-center justify-center font-bold ${
                  s < step ? 'bg-green-500 text-white' :
                  s === step ? 'bg-blue-600 text-white' :
                  'bg-gray-200 text-gray-600'
                }`}>
                  {s < step ? <Check className="w-6 h-6" /> : s}
                </div>
                <span className={`ml-3 font-medium ${s === step ? 'text-gray-900' : 'text-gray-500'}`}>
                  {['Project', 'Capabilities', 'Enablers', 'Resources', 'Review'][s - 1]}
                </span>
              </div>
              {s < 5 && <ChevronRight className="w-5 h-5 text-gray-400 mx-4" />}
            </React.Fragment>
          ))}
        </div>
      </div>

      {/* Step Content */}
      <div className="bg-white rounded-lg border border-gray-200 p-8">
        {step === 1 && <Step1ProjectDetails formData={formData} updateFormData={updateFormData} />}
        {step === 2 && <Step2CapabilitySelection formData={formData} toggleCapability={toggleCapability} />}
        {step === 3 && <Step3EnablerConfiguration formData={formData} />}
        {step === 4 && <Step4ResourceRequirements formData={formData} updateFormData={updateFormData} />}
        {step === 5 && <Step5Review formData={formData} calculateEstimatedCost={calculateEstimatedCost} />}

        {/* Navigation */}
        <div className="flex items-center justify-between mt-8 pt-6 border-t border-gray-200">
          <button
            onClick={() => step > 1 ? setStep(step - 1) : setView('dashboard')}
            className="flex items-center px-6 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition font-medium"
          >
            <ChevronLeft className="w-5 h-5 mr-1" />
            {step === 1 ? 'Cancel' : 'Previous'}
          </button>
          <div className="text-sm text-gray-600">
            Step {step} of {totalSteps}
          </div>
          {step < totalSteps ? (
            <button
              onClick={() => setStep(step + 1)}
              disabled={step === 1 && !formData.name}
              className="flex items-center px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
              <ChevronRight className="w-5 h-5 ml-1" />
            </button>
          ) : (
            <button
              onClick={handleSubmit}
              className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition font-medium"
            >
              Create Environment
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

// Step 1: Project Details
function Step1ProjectDetails({ formData, updateFormData }) {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Project Details</h2>
        <p className="text-gray-600">Basic information about your simulation environment</p>
      </div>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Environment Name *
          </label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => updateFormData('name', e.target.value)}
            placeholder="e.g., Production Load Test Environment"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            value={formData.description}
            onChange={(e) => updateFormData('description', e.target.value)}
            placeholder="Describe the purpose and goals of this environment..."
            rows={4}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Owner/Team
            </label>
            <input
              type="text"
              value={formData.owner}
              onChange={(e) => updateFormData('owner', e.target.value)}
              placeholder="Engineering Team"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Tags
            </label>
            <input
              type="text"
              value={formData.tags}
              onChange={(e) => updateFormData('tags', e.target.value)}
              placeholder="production, load-test, critical"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

// Step 2: Capability Selection
function Step2CapabilitySelection({ formData, toggleCapability }) {
  const [filter, setFilter] = useState('');

  const getDependencyStatus = (cap) => {
    if (cap.dependencies.length === 0) return 'ready';
    const allMet = cap.dependencies.every(depId => formData.capabilities.includes(depId));
    return allMet ? 'ready' : 'blocked';
  };

  const filteredCapabilities = CAPABILITIES.filter(cap =>
    cap.name.toLowerCase().includes(filter.toLowerCase()) ||
    cap.description.toLowerCase().includes(filter.toLowerCase()) ||
    cap.id.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Select Capabilities</h2>
        <p className="text-gray-600">Choose the capabilities your simulation environment needs</p>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
        <input
          type="text"
          placeholder="Search capabilities..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      <div className="text-sm text-gray-600 flex items-center space-x-4">
        <span className="font-medium">Selected: {formData.capabilities.length}</span>
        <span>•</span>
        <span>Required Enablers: {Object.keys(formData.enablers).length}</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-96 overflow-y-auto pr-2">
        {filteredCapabilities.map((cap) => {
          const isSelected = formData.capabilities.includes(cap.id);
          const depStatus = getDependencyStatus(cap);
          const isBlocked = depStatus === 'blocked' && !isSelected;

          return (
            <div
              key={cap.id}
              onClick={() => !isBlocked && toggleCapability(cap.id)}
              className={`p-4 border-2 rounded-lg cursor-pointer transition ${
                isSelected
                  ? 'border-blue-500 bg-blue-50'
                  : isBlocked
                  ? 'border-gray-200 bg-gray-50 opacity-60 cursor-not-allowed'
                  : 'border-gray-200 hover:border-blue-300 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center space-x-2">
                    <span className="text-xs font-mono text-gray-500">{cap.id}</span>
                    {isBlocked && (
                      <AlertCircle className="w-4 h-4 text-orange-500" />
                    )}
                  </div>
                  <h3 className="font-semibold text-gray-900 mt-1">{cap.name}</h3>
                  <p className="text-sm text-gray-600 mt-1">{cap.description}</p>
                  {cap.dependencies.length > 0 && (
                    <div className="mt-2 text-xs text-gray-500">
                      Requires: {cap.dependencies.join(', ')}
                    </div>
                  )}
                </div>
                <div className={`w-6 h-6 rounded border-2 flex items-center justify-center flex-shrink-0 ml-3 ${
                  isSelected ? 'border-blue-500 bg-blue-500' : 'border-gray-300'
                }`}>
                  {isSelected && <Check className="w-4 h-4 text-white" />}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// Step 3: Enabler Configuration
function Step3EnablerConfiguration({ formData }) {
  const selectedEnablers = Object.keys(formData.enablers).sort();

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Enabler Configuration</h2>
        <p className="text-gray-600">
          Auto-selected based on your capabilities ({selectedEnablers.length} enablers required)
        </p>
      </div>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex items-start">
          <AlertCircle className="w-5 h-5 text-blue-600 mt-0.5 mr-3 flex-shrink-0" />
          <div className="text-sm text-blue-800">
            These enablers are automatically selected based on your capability choices. They provide the foundational services and integrations needed for your environment.
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-96 overflow-y-auto pr-2">
        {selectedEnablers.map((enablerId) => {
          const enabler = ENABLERS[enablerId];
          return (
            <div key={enablerId} className="p-4 border-2 border-green-500 bg-green-50 rounded-lg">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center space-x-2">
                    <span className="text-xs font-mono text-gray-600">{enablerId}</span>
                    <span className="px-2 py-0.5 bg-green-200 text-green-800 text-xs rounded-full font-medium">
                      Required
                    </span>
                  </div>
                  <h3 className="font-semibold text-gray-900 mt-1">{enabler.name}</h3>
                  <p className="text-sm text-gray-600 mt-1">{enabler.description}</p>
                </div>
                <Check className="w-6 h-6 text-green-600 flex-shrink-0 ml-3" />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// Step 4: Resource Requirements
function Step4ResourceRequirements({ formData, updateFormData }) {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Resource Requirements</h2>
        <p className="text-gray-600">Define compute, storage, and networking needs</p>
      </div>

      <div className="space-y-6">
        {/* Compute */}
        <div className="p-6 border border-gray-200 rounded-lg">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Server className="w-5 h-5 mr-2 text-blue-600" />
            Compute Resources
          </h3>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">CPU Cores</label>
              <input
                type="number"
                value={formData.compute.cpu}
                onChange={(e) => updateFormData('compute', { ...formData.compute, cpu: parseInt(e.target.value) })}
                min="1"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Memory (GB)</label>
              <input
                type="number"
                value={formData.compute.memory}
                onChange={(e) => updateFormData('compute', { ...formData.compute, memory: parseInt(e.target.value) })}
                min="1"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Instances</label>
              <input
                type="number"
                value={formData.compute.instances}
                onChange={(e) => updateFormData('compute', { ...formData.compute, instances: parseInt(e.target.value) })}
                min="1"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        {/* Storage */}
        <div className="p-6 border border-gray-200 rounded-lg">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Database className="w-5 h-5 mr-2 text-purple-600" />
            Storage
          </h3>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Storage Size (GB)</label>
            <input
              type="number"
              value={formData.storage}
              onChange={(e) => updateFormData('storage', parseInt(e.target.value))}
              min="1"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        {/* Network */}
        <div className="p-6 border border-gray-200 rounded-lg">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Network className="w-5 h-5 mr-2 text-green-600" />
            Network Configuration
          </h3>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Network Type</label>
            <select
              value={formData.network}
              onChange={(e) => updateFormData('network', e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="private">Private Network</option>
              <option value="public">Public Network</option>
              <option value="hybrid">Hybrid Network</option>
            </select>
          </div>
        </div>

        {/* Scheduling */}
        <div className="p-6 border border-gray-200 rounded-lg">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Clock className="w-5 h-5 mr-2 text-orange-600" />
            Scheduling
          </h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Priority</label>
              <select
                value={formData.priority}
                onChange={(e) => updateFormData('priority', e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="critical">Critical</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Duration (hours)</label>
              <input
                type="number"
                value={formData.duration}
                onChange={(e) => updateFormData('duration', parseInt(e.target.value))}
                min="1"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Step 5: Review
function Step5Review({ formData, calculateEstimatedCost }) {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Review & Submit</h2>
        <p className="text-gray-600">Review your environment configuration before creation</p>
      </div>

      {/* Cost Estimate */}
      <div className="bg-gradient-to-r from-blue-50 to-purple-50 border-2 border-blue-200 rounded-lg p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-700">Estimated Daily Cost</p>
            <p className="text-4xl font-bold text-gray-900 mt-1">${calculateEstimatedCost()}</p>
            <p className="text-sm text-gray-600 mt-1">Based on selected resources and capabilities</p>
          </div>
          <DollarSign className="w-16 h-16 text-blue-600 opacity-50" />
        </div>
      </div>

      {/* Summary Sections */}
      <div className="space-y-4">
        {/* Project Info */}
        <div className="border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Project Information</h3>
          <dl className="grid grid-cols-2 gap-4">
            <div>
              <dt className="text-sm font-medium text-gray-500">Name</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.name || 'Not specified'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Owner</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.owner || 'Not specified'}</dd>
            </div>
            <div className="col-span-2">
              <dt className="text-sm font-medium text-gray-500">Description</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.description || 'No description'}</dd>
            </div>
          </dl>
        </div>

        {/* Capabilities */}
        <div className="border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Capabilities ({formData.capabilities.length})
          </h3>
          <div className="flex flex-wrap gap-2">
            {formData.capabilities.map(capId => {
              const cap = CAPABILITIES.find(c => c.id === capId);
              return (
                <span key={capId} className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium">
                  {cap.id}: {cap.name}
                </span>
              );
            })}
          </div>
        </div>

        {/* Resources */}
        <div className="border border-gray-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Resource Allocation</h3>
          <dl className="grid grid-cols-2 gap-4">
            <div>
              <dt className="text-sm font-medium text-gray-500">Compute</dt>
              <dd className="mt-1 text-sm text-gray-900">
                {formData.compute.cpu} CPU × {formData.compute.memory} GB RAM × {formData.compute.instances} instances
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Storage</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.storage} GB</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Network</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.network}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Priority</dt>
              <dd className="mt-1 text-sm text-gray-900">{formData.priority}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  );
}

// Templates Library
function TemplatesLibrary() {
  const templates = [
    {
      id: 'tmpl-001',
      name: 'Standard Load Test',
      description: 'Pre-configured environment for load and stress testing',
      capabilities: ['C01', 'C02', 'C03', 'C04', 'C06', 'C12'],
      popularity: 145,
      cost: 98.50
    },
    {
      id: 'tmpl-002',
      name: 'Full Production Mirror',
      description: 'Complete production-like environment with all capabilities',
      capabilities: ['C01', 'C02', 'C03', 'C04', 'C05', 'C06', 'C07', 'C08', 'C09', 'C12'],
      popularity: 89,
      cost: 287.20
    },
    {
      id: 'tmpl-003',
      name: 'Dev/Test Basic',
      description: 'Minimal environment for development and basic testing',
      capabilities: ['C01', 'C02', 'C04', 'C12'],
      popularity: 203,
      cost: 45.00
    }
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Template Library</h2>
          <p className="text-gray-600 mt-1">Pre-configured environments to get started quickly</p>
        </div>
        <div className="relative">
          <Search className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Search templates..."
            className="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {templates.map((template) => (
          <div key={template.id} className="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg transition">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-gray-900">{template.name}</h3>
                <p className="text-sm text-gray-600 mt-2">{template.description}</p>
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600">Capabilities:</span>
                <span className="font-medium text-gray-900">{template.capabilities.length}</span>
              </div>
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600">Used by:</span>
                <span className="font-medium text-gray-900">{template.popularity} teams</span>
              </div>
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600">Est. cost/day:</span>
                <span className="font-medium text-green-600">${template.cost}</span>
              </div>
            </div>

            <button className="w-full mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-medium">
              Use Template
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
