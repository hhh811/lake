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
import { useCallback, useEffect, useState } from 'react'
import parser from 'cron-parser'
import { BlueprintMode } from '@/data/NullBlueprint'
import { Providers } from '@/data/Providers'

function useBlueprintValidation ({
  name,
  cronConfig,
  customCronConfig,
  enable,
  tasks = [],
  mode = null,
  connections = [],
  entities = {},
  boards = {},
  projects = {},
  activeStep = null,
  activeProvider = null,
  activeConnection = null
}) {
  const [errors, setErrors] = useState([])
  const [isValid, setIsValid] = useState(false)

  const clear = () => {
    setErrors([])
  }

  const isValidCronExpression = useCallback((expression) => {
    let isValid = false
    try {
      parser.parseExpression(expression)
      isValid = true
    } catch (e) {
      isValid = false
    }
    return isValid
  }, [])

  const validateNumericSet = useCallback((set = []) => {
    return Array.isArray(set) ? set.every(i => !isNaN(i)) : false
  }, [])

  const validateRepositoryName = useCallback((set = []) => {
    const repoRegExp = /([a-z0-9_-]){2,}\/([a-z0-9_-]){2,}$/gi
    return set.every(i => i.match(repoRegExp))
  }, [])

  const valiateNonEmptySet = useCallback((set = []) => {
    return set.length > 0
  }, [])

  const validate = useCallback(() => {
    const errs = []

    if (!name) {
      errs.push('Blueprint Name: Enter a valid Name')
    }

    if (name && name.length <= 2) {
      errs.push('Blueprint Name: Name too short, 3 chars minimum.')
    }

    if (mode !== null && ![BlueprintMode.NORMAL, BlueprintMode.ADVANCED].includes(mode)) {
      errs.push('Invalid / Unsupported Blueprint Mode Detected!')
    }

    if (mode === BlueprintMode.NORMAL) {
      if (!cronConfig) {
        errs.push('Blueprint Cron: No Crontab schedule defined.')
      }

      if (cronConfig && !['custom', 'manual'].includes(cronConfig) && !isValidCronExpression(cronConfig)) {
        errs.push('Blueprint Cron: Invalid Crontab Expression, unable to parse.')
      }

      if (cronConfig === 'custom' && !isValidCronExpression(customCronConfig)) {
        errs.push(`Blueprint Cron: Invalid Custom Expression, unable to parse. [${customCronConfig}]`)
      }

      if (enable && tasks?.length === 0) {
        errs.push('Blueprint Tasks: Invalid/Empty Configuration')
      }

      switch (activeStep?.id) {
        case 1:
          if (connections.length === 0) {
            errs.push('No Data Connections selected.')
          }
          break
        case 2:
          if (activeProvider?.id === Providers.JIRA && boards[activeConnection?.id]?.length === 0) {
            errs.push('Boards: No Boards selected.')
          }
          if (activeProvider?.id === Providers.GITHUB && projects[activeConnection?.id]?.length === 0) {
            errs.push('Projects: No Project Repsitories entered.')
          }
          if (activeProvider?.id === Providers.GITHUB && !validateRepositoryName(projects[activeConnection?.id])) {
            errs.push('Projects: Only Git Repository Names are supported (username/repo).')
          }
          if (entities[activeConnection?.id]?.length === 0) {
            errs.push('Data Entities: No Data Entities selected.')
          }
          if (activeProvider?.id === Providers.GITLAB && projects[activeConnection?.id]?.length === 0) {
            errs.push('Projects: No Project IDs entered.')
          }
          if (activeProvider?.id === Providers.GITLAB && !validateNumericSet(projects[activeConnection?.id])) {
            errs.push('Projects: Only Numeric Project IDs are supported.')
          }

          connections.forEach(c => {
            if (c.provider === Providers.JIRA && boards[c?.id]?.length === 0) {
              errs.push(`${c.name} requires a Board`)
            }
            if (c.provider === Providers.GITHUB && projects[c?.id]?.length === 0) {
              errs.push(`${c.name} requires Project Names`)
            }
            if (c.provider === Providers.GITHUB && !validateRepositoryName(projects[c?.id])) {
              errs.push(`${c.name} has Invalid Project Repository`)
            }
            if (c.provider === Providers.GITLAB && projects[c?.id]?.length === 0) {
              errs.push(`${c.name} requires Project IDs`)
            }
            if (c.provider === Providers.GITLAB && !validateNumericSet(projects[c?.id])) {
              errs.push(`${c.name} has invalid Project ID`)
            }
            if (entities[c?.id]?.length === 0) {
              errs.push(`${c.name} is missing Data Entities`)
            }
          })

          break
      }
    }

    setErrors(errs)
  }, [
    name,
    cronConfig,
    customCronConfig,
    tasks,
    enable,
    mode,
    connections,
    boards,
    entities,
    projects,
    activeStep,
    activeProvider?.id,
    activeConnection,
    isValidCronExpression,
    validateNumericSet,
    validateRepositoryName
  ])

  const fieldHasError = useCallback((fieldId) => {
    return errors.some(e => e.includes(fieldId))
  }, [errors])

  const getFieldError = useCallback((fieldId) => {
    return errors.find(e => e.includes(fieldId))
  }, [errors])

  useEffect(() => {
    // console.log('>>> BLUEPRINT FORM ERRORS...', errors)
    setIsValid(errors.length === 0)
    if (errors.length > 0) {
      // ToastNotification.clear()
    }
  }, [errors])

  return {
    errors,
    setErrors,
    isValid,
    validate,
    clear,
    fieldHasError,
    getFieldError
  }
}

export default useBlueprintValidation
