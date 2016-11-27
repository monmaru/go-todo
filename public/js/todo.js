(function(Vue) {
  "use strict";
  
  new Vue({
    el: 'body',

    data: {
      todos: [],
      newTodo: {}
    },

    created: function() {
      this.getAllTodos();
    },

    methods: {
      getAllTodos: function() {
        this.$http.get('/api/todos').then(function(res) {
          this.todos = res.data ? res.data : [];
        });
      },

      createTodo: function() {
        if (!$.trim(this.newTodo.text)) {
          this.newTodo = {};
          return;
        };

        this.newTodo.done = false;

        this.$http.post('/api/todos', this.newTodo).then(function(res) {
          this.newTodo.id = res.data.id;
          this.todos.push(this.newTodo);
          this.newTodo = {};
        }).catch(function(err) {
          console.log(err);
        });
      },

      deleteTodo: function(id) {
        this.$http.delete('/api/todos/' + id).then(function(res) {
          this.getAllTodos();
        }).catch(function(err) {
          console.log(err);
        });
      },

      updateTodo: function(todo, completed) {
         todo.done = completed;
         this.$http.put('/api/todos', todo).then(function(res) {
           this.getAllTodos();
         }).catch(function(err) {
           console.log(err);
        });
      }
    }
  });
})(Vue);
