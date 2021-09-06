class BaseController {
  constructor(instance) {
    this.instance = instance;
  }
  findAll(attributes) {
    return this.instance.findAll(attributes);
  }
  findByFilter(attributes, where) {
    return this.instance.findByFilter(attributes, where);
  }
  findByFilterOrder(attributes, where, order) {
    return this.instance.findByFilterOrder(attributes, where, order);
  }
  findLikeByFilter(attributes, where) {
    return this.instance.findLikeByFilter(attributes, where);
  }
  findLikeByFilterOrder(attributes, where, order) {
    return this.instance.findLikeByFilterOrder(attributes, where, order);
  }
  update(attributes, where) {
    return this.instance.update(attributes, where);
  }
  delete(where) {
    return this.instance.delete(where);
  }
  create(entity) {
    return this.instance.create(entity);
  }
  createBatch(entities) {
    return this.instance.createBatch(entities);
  }
}
module.exports = BaseController;
