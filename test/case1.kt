TypeDecl:
    class A {
    }

expect:
    CLASS(
      IDENT('A'),
      MEMBERS
    )
