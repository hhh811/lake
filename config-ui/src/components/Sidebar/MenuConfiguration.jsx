/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import { ProviderLabels } from '@/data/Providers'
import { GRAFANA_URL } from '@/utils/config'

const MenuConfiguration = (activeRoute) => {
  return [
    {
      id: 0,
      label: 'Data Integrations',
      route: '/integrations',
      active: activeRoute.url.startsWith('/integrations') || activeRoute.url === '/',
      icon: 'data-connection',
      classNames: [],
      children: [
        {
          id: 0,
          label: ProviderLabels.JIRA,
          route: '/integrations/jira',
          active: activeRoute.url.endsWith('/integrations/jira') || activeRoute.url.endsWith('/jira'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 1,
          label: ProviderLabels.GITHUB,
          route: '/integrations/github',
          active: activeRoute.url.endsWith('/integrations/github') || activeRoute.url.endsWith('/github'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 2,
          label: ProviderLabels.GITLAB,
          route: '/integrations/gitlab',
          active: activeRoute.url.endsWith('/integrations/gitlab') || activeRoute.url.endsWith('/gitlab'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 3,
          label: ProviderLabels.JENKINS,
          route: '/integrations/jenkins',
          active: activeRoute.url.endsWith('/integrations/jenkins') || activeRoute.url.endsWith('/jenkins'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 4,
          label: ProviderLabels.TAPD,
          route: '/integrations/tapd',
          active: activeRoute.url.endsWith('/integrations/tapd') || activeRoute.url.endsWith('/tapd'),
          icon: 'layers',
          classNames: [],
        }
      ]
    },
    {
      id: 1,
      label: 'Blueprints',
      icon: 'home',
      route: '/blueprints',
      disabled: false,
      active: activeRoute.url === '/blueprints',
      children: [
        {
          id: 0,
          label: 'Create Bluepint',
          route: '/blueprints/create',
          active: activeRoute.url.endsWith('/blueprints/create'),
          icon: 'git-pull',
          classNames: [],
        },
      ]
    },
    {
      id: 2,
      label: 'Connections',
      disabled: true,
      icon: 'git-merge',
      classNames: [],
      route: '/connections',
      active: activeRoute.url === '/connections',
      children: [
      ]
    },
    {
      id: 2,
      label: 'Dashboard',
      icon: 'dashboard',
      classNames: [],
      route: GRAFANA_URL,
      active: false,
      children: [
      ]
    },
    // {
    //   id: 3,
    //   label: 'Pipelines',
    //   icon: 'git-merge',
    //   classNames: [],
    //   route: '/pipelines',
    //   active: activeRoute.url.startsWith('/pipelines'),
    //   children: [
    //     {
    //       id: 0,
    //       label: 'Create Pipeline Run',
    //       route: '/pipelines/create',
    //       active: activeRoute.url.endsWith('/pipelines/create'),
    //       icon: 'git-pull',
    //       classNames: [],
    //     },
    //     {
    //       id: 1,
    //       label: 'All Pipeline Runs',
    //       route: '/pipelines',
    //       active: activeRoute.url.endsWith('/pipelines'),
    //       icon: 'layers',
    //       classNames: [],
    //       disabled: false
    //     },
    //     {
    //       id: 2,
    //       label: 'Pipeline Blueprints',
    //       route: '/blueprints',
    //       active: activeRoute.url.endsWith('/blueprints'),
    //       icon: 'grid',
    //       classNames: [],
    //     },
    //   ]
    // },
    // {
    //   id: 3,
    //   label: 'Documentation',
    //   icon: 'help',
    //   classNames: [],
    //   route: 'https://github.com/apache/incubator-devlake/wiki',
    //   target: "_blank",
    //   external: true,
    //   active: activeRoute.url === '/documentation',
    //   children: [
    //   ]
    // },
  ]
}

export default MenuConfiguration
