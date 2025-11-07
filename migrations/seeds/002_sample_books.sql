-- Seed sample books data
-- Created at: 2024-10-25

INSERT INTO books (title, author, isbn, year, pages, publisher, description)
VALUES
    (
        'The Go Programming Language',
        'Alan Donovan and Brian Kernighan',
        '978-0134190440',
        2015,
        400,
        'Addison-Wesley Professional',
        'The Go Programming Language is the authoritative resource for any programmer who wants to learn Go.'
    ),
    (
        'Clean Code: A Handbook of Agile Software Craftsmanship',
        'Robert C. Martin',
        '978-0132350884',
        2008,
        464,
        'Prentice Hall',
        'Even bad code can function. But if code isn''t clean, it can bring a development organization to its knees.'
    ),
    (
        'Design Patterns: Elements of Reusable Object-Oriented Software',
        'Gang of Four',
        '978-0201633610',
        1994,
        416,
        'Addison-Wesley',
        'Design Patterns is a modern classic in the literature of object-oriented development.'
    ),
    (
        'Refactoring: Improving the Design of Existing Code',
        'Martin Fowler',
        '978-0201485677',
        1999,
        464,
        'Addison-Wesley Longman Publishing Co.',
        'As the application of object technology--particularly the Java programming language--has become commonplace, a new problem has emerged to trouble the Enterprise.'
    ),
    (
        'The Pragmatic Programmer: Your Journey to Mastery',
        'David Thomas and Andrew Hunt',
        '978-0201616224',
        1999,
        352,
        'Addison-Wesley',
        'The Pragmatic Programmer cuts through the increasing specialization and technicalities of modern software development.'
    )
ON CONFLICT (isbn) DO NOTHING;
