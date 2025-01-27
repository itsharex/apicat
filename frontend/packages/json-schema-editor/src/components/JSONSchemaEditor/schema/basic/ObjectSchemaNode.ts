import { JSONSchema } from '../../types'
import BasicTypeSchemaNode from '../BasicTypeSchemaNode'
import SchemaNode from '../SchemaNode'

export default class ObjectSchemaNode extends BasicTypeSchemaNode {
  type = 'object'

  createDefaultSchema(): JSONSchema {
    const schema = super.createDefaultSchema()
    schema.properties = {}
    schema.required = []
    schema['x-apicat-orders'] = []
    console.log(schema)
    return schema
  }

  removeChild(child: SchemaNode): void {
    const index = this.childNodes.indexOf(child)

    if (index > -1) {
      this.store && this.store.deregisterNode(child)
      child.parent = null
      this.childNodes.splice(index, 1)
      this.changeNotify()
    }
  }

  // 删除属性
  deleteProperty() {
    const { schema } = this
    if (!schema.properties) {
      return
    }
    // remove the property
    delete schema.properties[this.name]
    // remove the required
    schema.required = schema.required?.filter((one) => one !== this.name)
    // remove the x-apicat-orders
    schema['x-apicat-orders'] = schema['x-apicat-orders']?.filter((one) => one !== this.name)
  }

  // 增加属性
  addProperty(child: SchemaNode) {
    if (!child) {
      return
    }
    const property = this.schema.properties || {}
    property[child.name] = child.schema
    this.schema.properties = property
  }

  // 更新属性
  updatePropertyName(newChild: SchemaNode, oldChild: SchemaNode) {
    // const property = this.schema.properties || {}
    // property[child.name] = child.schema
    // this.schema.properties = property
  }
}
