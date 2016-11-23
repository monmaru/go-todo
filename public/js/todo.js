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

        this.$http.post('/api/todos', this.newTodo).success(function(res) {
          this.newTodo.id = res.id;
          this.todos.push(this.newTodo);

          this.newTodo = {};
        }).error(function(err) {
          console.log(err);
        });
      },

      deleteTodo: function(index) {
        this.$http.delete('/api/todos/' + index).success(function(res) {
          this.getAllTodos();
        }).error(function(err) {
          console.log(err);
        });
      },

      updateTodo: function(todo, completed) {
        todo.done = completed;
        this.$http.post('/api/todos', todo).success(function(res) {
          this.getAllTodos();
        }).error(function(err) {
          console.log(err);
        });
      }
    }
  });
})(Vue);
