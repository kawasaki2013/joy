;(function() {
  var dep2 = (function() {
    // AnotherGet
    function AnotherGet() {
      return 'another world'
    }

    return {
      AnotherGet: AnotherGet
    }
  })()

  var dep = (function() {
    // Get fn
    function Get() {
      return 'world'
    }

    return {
      Get: Get
    }
  })()

  var main = (function() {
    function main() {
      console.log('hi ' + dep.Get() + ' ' + Side() + ' ' + dep2.AnotherGet())
    }

    function Side() {
      return 'Side'
    }

    return {
      main: main,
      Side: Side
    }
  })()

  main.main()
})()
