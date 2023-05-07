import useApi from '@/hooks/useApi'
import useCommonResponseStore from '@/store/commonResponse'
import { ProjectInfo, APICatCommonResponseCustom, APICatCommonResponse } from '@/typings'
import { createCommonResponse } from '@/views/document/components/createHttpDocument'
import { uuid } from '@apicat/shared'
import { useI18n } from 'vue-i18n'

export const useResponseparamList = ({ id: project_id }: Pick<ProjectInfo, 'id'>) => {
  const { t } = useI18n()
  const commonResponseStore = useCommonResponseStore()
  const [isLoading, getResponseParamListApi] = useApi(commonResponseStore.getCommonResponseList)()

  const responseParamList: Ref<APICatCommonResponseCustom[]> = ref([])
  const isLoadingForDelete = shallowRef(false)

  const extendResponseParamModel = (param?: Partial<APICatCommonResponseCustom>): APICatCommonResponseCustom => {
    return {
      id: param?.id ?? uuid(),
      isLocal: true,
      expand: false,
      isLoaded: false,
      isLoading: false,
      ...param,
    }
  }

  const createResponseParamModel = () => {
    const response = createCommonResponse({ name: t('app.response.model.name') })
    const extendModel = extendResponseParamModel({ code: response.code, description: response.description, isLoaded: true, expand: true })
    extendModel.detail = response
    return extendModel
  }

  const handleAddParam = () => {
    responseParamList.value.unshift(createResponseParamModel())
  }

  const handleDeleteParam = async (item: APICatCommonResponseCustom, index: number, isUnRef: number) => {
    const { isLocal } = item

    if (!isLocal) {
      isLoadingForDelete.value = true
      try {
        await commonResponseStore.deleteResponseParam(project_id, { id: item.id } as any, isUnRef)
      } finally {
        isLoadingForDelete.value = false
      }
    }

    responseParamList.value.splice(index, 1)
  }

  onMounted(async () => {
    const data: APICatCommonResponse[] = await getResponseParamListApi(project_id)
    const list: APICatCommonResponseCustom[] = data.map((item) => extendResponseParamModel({ ...item, isLocal: false }))
    responseParamList.value = list
  })

  return {
    isLoading,
    isLoadingForDelete,
    getResponseParamListApi,
    responseParamList,

    handleAddParam,
    handleDeleteParam,
  }
}