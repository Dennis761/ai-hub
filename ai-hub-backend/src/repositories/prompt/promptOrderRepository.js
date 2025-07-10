import PromptModel from '../../models/promptModel.js';

class PromptOrderRepository {
  bulkUpdateExecutionOrder(orderArray) {
    const bulkOps = orderArray.map((item, index) => ({
      updateOne: {
        filter: { _id: item._id },
        update: { executionOrder: index + 1 }
      }
    }));

    return PromptModel.bulkWrite(bulkOps);
  }
}

export default PromptOrderRepository;
