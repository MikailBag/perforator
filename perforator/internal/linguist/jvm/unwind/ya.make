IF (NOT CI)
    RECURSE(sample)
ENDIF()

RECURSE(
    cheatsheets/tool
    configure
    lib
)

ENDIF()
